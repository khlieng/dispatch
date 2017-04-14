// Copyright 2015, Joe Tsai. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE.md file.

package brotli

import "io"
import "io/ioutil"
import "bufio"
import "bytes"
import "strings"
import "encoding/hex"
import "runtime"
import "testing"

func TestReader(t *testing.T) {
	var vectors = []struct {
		desc   string // Description of the test
		input  string // Test input string in hex
		output string // Expected output string in hex
		inIdx  int64  // Expected input offset after reading
		outIdx int64  // Expected output offset after reading
		err    error  // Expected error
	}{{
		desc: "empty string (truncated)",
		err:  io.ErrUnexpectedEOF,
	}, {
		desc:  "empty last block (WBITS: 16)",
		input: "06",
		inIdx: 1,
	}, {
		desc:  "empty last block (WBITS: 12)",
		input: "c101",
		inIdx: 2,
	}, {
		desc:  "empty last block (WBITS: 17)",
		input: "8101",
		inIdx: 2,
	}, {
		desc:  "empty last block (WBITS: 21)",
		input: "39",
		inIdx: 1,
	}, {
		desc:  "empty last block (WBITS: invalid)",
		input: "9101",
		inIdx: 1,
		err:   ErrCorrupt,
	}, {
		desc:  "empty last block (trash at the end)",
		input: "06ff",
		inIdx: 1,
	}, {
		desc:  "empty last block (padding is non-zero)",
		input: "16",
		inIdx: 1,
		err:   ErrCorrupt,
	}, {
		desc:  "empty meta data block (MLEN: 0)",
		input: "0c03",
		inIdx: 2,
	}, {
		desc:  "meta data block",
		input: "2c0648656c6c6f2c20776f726c642103",
		inIdx: 16,
	}, {
		desc:  "meta data block (truncated)",
		input: "2c06",
		inIdx: 2,
		err:   io.ErrUnexpectedEOF,
	}, {
		desc:  "meta data block (use reserved bit)",
		input: "3c0648656c6c6f2c20776f726c642103",
		inIdx: 1,
		err:   ErrCorrupt,
	}, {
		desc:  "meta data block (meta padding is non-zero)",
		input: "2c8648656c6c6f2c20776f726c642103",
		inIdx: 2,
		err:   ErrCorrupt,
	}, {
		desc:  "meta data block (non-minimal MLEN)",
		input: "4c060048656c6c6f2c20776f726c642103",
		inIdx: 3,
		err:   ErrCorrupt,
	}, {
		desc:  "meta data block (MLEN: 1<<0)",
		input: "2c00ff03",
		inIdx: 4,
	}, {
		desc:  "meta data block (MLEN: 1<<24)",
		input: "ecffff7f" + strings.Repeat("f0", 1<<24) + "03",
		inIdx: 5 + 1<<24,
	}, {
		desc:   "raw data block",
		input:  "c0001048656c6c6f2c20776f726c642103",
		output: "48656c6c6f2c20776f726c6421",
		inIdx:  17,
		outIdx: 13,
	}, {
		desc:  "raw data block (truncated)",
		input: "c00010",
		inIdx: 3,
		err:   io.ErrUnexpectedEOF,
	}, {
		desc:  "raw data block (raw padding is non-zero)",
		input: "c000f048656c6c6f2c20776f726c642103",
		inIdx: 3,
		err:   ErrCorrupt,
	}, {
		desc:  "raw data block (non-minimal MLEN)",
		input: "c400000148656c6c6f2c20776f726c642103",
		inIdx: 3,
		err:   ErrCorrupt,
	}, {
		desc:   "raw data block (MLEN: 1<<0)",
		input:  "0000106103",
		output: "61",
		inIdx:  4 + 1<<0,
		outIdx: 1 << 0,
	}, {
		desc:   "raw data block (MLEN: 1<<24)",
		input:  "f8ffff1f" + strings.Repeat("f0", 1<<24) + "03",
		output: strings.Repeat("f0", 1<<24),
		inIdx:  5 + 1<<24,
		outIdx: 1 << 24,
	}, {
		desc:   "simple prefix (|L|:1 |I|:1 |D|:1 MLEN:1)",
		input:  "00000000c4682010c0",
		output: "a3",
		inIdx:  9,
		outIdx: 1,
	}, {
		desc:   "simple prefix, out-of-order (|L|:2 |I|:1 |D|:1 MLEN:1)",
		input:  "00000000d4a8682010c001",
		output: "a3",
		inIdx:  11,
		outIdx: 1,
	}, {
		desc:   "simple prefix, non-unique (|L|:2 |I|:1 |D|:1 MLEN:1)",
		input:  "00000000d4e8682010c001",
		output: "",
		inIdx:  7,
		outIdx: 0,
		err:    ErrCorrupt,
	}, {
		desc:   "simple prefix, out-of-order (|L|:3 |I|:1 |D|:1 MLEN:1)",
		input:  "0000000024e8e96820104003",
		output: "a3",
		inIdx:  12,
		outIdx: 1,
	}, {
		desc:   "simple prefix, out-of-order, no-tree-select (|L|:4 |I|:1 |D|:1 MLEN:1)",
		input:  "0000000034e8e968a840208006",
		output: "a3",
		inIdx:  13,
		outIdx: 1,
	}, {
		desc:   "simple prefix, out-of-order, yes-tree-select (|L|:4 |I|:1 |D|:1 MLEN:1)",
		input:  "0000000034e8e968e94020800d",
		output: "a3",
		inIdx:  13,
		outIdx: 1,
	}, {
		desc:   "simple prefix, max-sym-ok (|L|:1 |I|:2 |D|:1 MLEN:1)",
		input:  "00000000c46821f06b0006",
		output: "a3",
		inIdx:  11,
		outIdx: 1,
	}, {
		desc:   "simple prefix, max-sym-bad (|L|:1 |I|:2 |D|:1 MLEN:1)",
		input:  "00000000c46821006c0006",
		output: "",
		inIdx:  9,
		outIdx: 0,
		err:    ErrCorrupt,
	}, {
		desc:   "complex prefix, skip-zero, terminate-clens-codes (|L|:1 |I|:2 |D|:1 MLEN:1)",
		input:  "0000000070472010c001",
		output: "01",
		inIdx:  10,
		outIdx: 1,
	}, {
		desc:   "complex prefix, skip-zero, terminate-clens-codes (|L|:1 |I|:2 |D|:1 MLEN:1)",
		input:  "0000000070c01d080470",
		output: "01",
		inIdx:  10,
		outIdx: 1,
	}, {
		desc:   "complex prefix, skip-zero, terminate-clens-codes (|L|:1 |I|:2 |D|:1 MLEN:2)",
		input:  "1000000070c01d1004d0",
		output: "0100",
		inIdx:  10,
		outIdx: 2,
	}, {
		desc:   "complex prefix, skip-zero, terminate-codes (|L|:1 |I|:4 |D|:1 MLEN:3)",
		input:  "20000000b0c100000056151804700e",
		output: "030201",
		inIdx:  15,
		outIdx: 3,
	}, {
		desc:   "complex prefix, skip-zero, under-subscribed (|L|:1 |I|:4 |D|:1 MLEN:3)",
		input:  "20000000b0c1000000ae2a3008e01c",
		output: "",
		inIdx:  10,
		outIdx: 0,
		err:    ErrCorrupt,
	}, {
		desc:   "complex prefix, skip-zero, over-subscribed (|L|:1 |I|:4 |D|:1 MLEN:3)",
		input:  "20000000b0c1000000ac0a0c023807",
		output: "",
		inIdx:  10,
		outIdx: 0,
		err:    ErrCorrupt,
	}, {
		desc:   "complex prefix, skip-zero, single clens (|L|:1 |I|:256 |D|:1 MLEN:4)",
		input:  "30000000000000020001420000a5ff5503",
		output: "00a5ffaa",
		inIdx:  17,
		outIdx: 4,
	}, {
		desc:   "complex prefix, skip-zero, single clens (|L|:1 |I|:32 |D|:1 MLEN:4)",
		input:  "3000000000c001000004080100faf7",
		output: "00051f1b",
		inIdx:  15,
		outIdx: 4,
	}, {
		desc:   "complex prefix, skip-zero, single clens, zero clen (|L|:1 |I|:? |D|:1 MLEN:4)",
		input:  "30000000007000000004080100faf7",
		output: "",
		inIdx:  10,
		outIdx: 0,
		err:    ErrCorrupt,
	}, {
		desc:   "complex prefix, skip-zero, empty clens (|L|:1 |I|:? |D|:1 MLEN:4)",
		input:  "30000000000000000001420080fe3d",
		output: "",
		inIdx:  9,
		outIdx: 0,
		err:    ErrCorrupt,
	}, {
		desc:   "complex prefix, skip-zero, single clens, rep-last clen (|L|:1 |I|:256 |D|:1 MLEN:4)",
		input:  "3000000000002000006a014200aa33cc5503",
		output: "55cc33aa",
		inIdx:  18,
		outIdx: 4,
	}, {
		desc:   "complex prefix, skip-zero, single clens, rep-last clen, over-subscribed (|L|:1 |I|:257 |D|:1 MLEN:4)",
		input:  "300000000000200000aa014200aa33cc5503",
		output: "",
		inIdx:  10,
		outIdx: 0,
		err:    ErrCorrupt,
	}, {
		desc:   "complex prefix, skip-zero, single clens, rep-last clen, integer overflow (|L|:1 |I|:1018 |D|:1 MLEN:4)",
		input:  "3000000000002000002a070801a8ce30570d",
		output: "",
		inIdx:  11,
		outIdx: 0,
		err:    ErrCorrupt,
	}, {
		desc:   "complex prefix, skip-two, single clens, rep-last clen (|L|:1 |I|:256 |D|:1 MLEN:4)",
		input:  "3000000008000f00805a801080ea0c73d5",
		output: "55cc33aa",
		inIdx:  17,
		outIdx: 4,
	}, {
		desc:   "complex prefix, skip-three, single clens, rep-last clen (|L|:1 |I|:256 |D|:1 MLEN:4)",
		input:  "300000000cc00300a0162004a03ac35c35",
		output: "55cc33aa",
		inIdx:  17,
		outIdx: 4,
	}, {
		desc: "complex prefix, skip-zero, linear clens (|L|:1 |I|:16 |D|:1 MLEN:16)",
		input: "f000000050555555ffff8bd5169058d43cb2fadcf77f201480dabdeff7f7efbf" +
			"fffddffffbfffe7fffff01",
		output: "6162636465666768696a6b6c6d6e6f70",
		inIdx:  43,
		outIdx: 16,
	}, {
		desc: "complex prefix, skip-zero, mixed clens (|L|:1 |I|:192 |D|:1 MLEN:16)",
		input: "f000000050555555ffffe37a310f369a4d4b80756cc779b0619a02a1002c29ab" +
			"ec066084eee99dfd67d8ac18",
		output: "000240525356575e717a8abcbdbed7d9",
		inIdx:  44,
		outIdx: 16,
	}, {
		desc:   "compressed string: \"Hello, world! Hello, world!\"",
		input:  "1b1a00008c946ed6540dc2825426d942de6a9668ea996c961e00",
		output: "48656c6c6f2c20776f726c64212048656c6c6f2c20776f726c6421",
		inIdx:  26,
		outIdx: 27,
	}, {
		desc:   "compressed string (padding is non-zero): \"Hello, world! Hello, world!\"",
		input:  "1b1a00008c946ed6540dc2825426d942de6a9668ea996c961e80",
		output: "48656c6c6f2c20776f726c64212048656c6c6f2c20776f726c6421",
		inIdx:  26,
		outIdx: 27,
		err:    ErrCorrupt,
	}}

	for i, v := range vectors {
		input, _ := hex.DecodeString(v.input)
		rd, err := NewReader(bytes.NewReader(input), nil)
		if err != nil {
			t.Errorf("test %d, unexpected NewReader error: %v", i, err)
		}
		data, err := ioutil.ReadAll(rd)
		output := hex.EncodeToString(data)

		if err != v.err {
			t.Errorf("test %d, %s\nerror mismatch: got %v, want %v", i, v.desc, err, v.err)
		}
		if output != v.output {
			t.Errorf("test %d, %s\noutput mismatch:\ngot  %v\nwant %v", i, v.desc, output, v.output)
		}
		if rd.InputOffset != v.inIdx {
			t.Errorf("test %d, %s\ninput offset mismatch: got %d, want %d", i, v.desc, rd.InputOffset, v.inIdx)
		}
		if rd.OutputOffset != v.outIdx {
			t.Errorf("test %d, %s\noutput offset mismatch: got %d, want %d", i, v.desc, rd.OutputOffset, v.outIdx)
		}
	}
}

func TestReaderGolden(t *testing.T) {
	var vectors = []struct {
		input  string // Input filename
		output string // Output filename
	}{
		{"empty.br", "empty"},
		{"empty.00.br", "empty"},
		{"empty.01.br", "empty"},
		{"empty.02.br", "empty"},
		{"empty.03.br", "empty"},
		{"empty.04.br", "empty"},
		{"empty.05.br", "empty"},
		{"empty.06.br", "empty"},
		{"empty.07.br", "empty"},
		{"empty.08.br", "empty"},
		{"empty.09.br", "empty"},
		{"empty.10.br", "empty"},
		{"empty.11.br", "empty"},
		{"empty.12.br", "empty"},
		{"empty.13.br", "empty"},
		{"empty.14.br", "empty"},
		{"empty.15.br", "empty"},
		{"empty.16.br", "empty"},
		{"empty.17.br", "empty"},
		{"empty.18.br", "empty"},
		{"zeros.br", "zeros"},
		{"x.br", "x"},
		{"x.00.br", "x"},
		{"x.01.br", "x"},
		{"x.02.br", "x"},
		{"x.03.br", "x"},
		{"xyzzy.br", "xyzzy"},
		{"10x10y.br", "10x10y"},
		{"64x.br", "64x"},
		{"backward65536.br", "backward65536"},
		{"quickfox.br", "quickfox"},
		{"quickfox_repeated.br", "quickfox_repeated"},
		{"ukkonooa.br", "ukkonooa"},
		{"monkey.br", "monkey"},
		{"random_org_10k.bin.br", "random_org_10k.bin"},
		{"asyoulik.txt.br", "asyoulik.txt"},
		{"compressed_file.br", "compressed_file"},
		{"compressed_repeated.br", "compressed_repeated"},
		{"alice29.txt.br", "alice29.txt"},
		{"lcet10.txt.br", "lcet10.txt"},
		{"mapsdatazrh.br", "mapsdatazrh"},
		{"plrabn12.txt.br", "plrabn12.txt"},
	}

	for i, v := range vectors {
		input, err := ioutil.ReadFile("testdata/" + v.input)
		if err != nil {
			t.Errorf("test %d: %s\n%v", i, v.input, err)
			continue
		}
		output, err := ioutil.ReadFile("testdata/" + v.output)
		if err != nil {
			t.Errorf("test %d: %s\n%v", i, v.output, err)
			continue
		}

		rd, err := NewReader(bytes.NewReader(input), nil)
		if err != nil {
			t.Errorf("test %d, unexpected NewReader error: %v", i, err)
		}
		data, err := ioutil.ReadAll(rd)
		if err != nil {
			t.Errorf("test %d, %s\nerror mismatch: got %v, want nil", i, v.input, err)
		}
		if string(data) != string(output) {
			t.Errorf("test %d, %s\noutput mismatch:\ngot  %q\nwant %q", i, v.input, string(data), string(output))
		}
	}
}

func benchmarkDecode(b *testing.B, testfile string) {
	b.StopTimer()
	b.ReportAllocs()

	input, err := ioutil.ReadFile("testdata/" + testfile)
	if err != nil {
		b.Fatal(err)
	}
	r, err := NewReader(bytes.NewReader(input), nil)
	if err != nil {
		b.Fatal(err)
	}
	output, err := ioutil.ReadAll(r)
	if err != nil {
		b.Fatal(err)
	}

	nb := int64(len(output))
	output = nil
	runtime.GC()

	b.SetBytes(nb)
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		r, err := NewReader(bufio.NewReader(bytes.NewReader(input)), nil)
		if err != nil {
			b.Fatal(err)
		}
		cnt, err := io.Copy(ioutil.Discard, r)
		if err != nil {
			b.Fatalf("unexpected error: %v", err)
		}
		if cnt != nb {
			b.Fatalf("unexpected count: got %d, want %d", cnt, nb)
		}
	}
}

func BenchmarkDecodeDigitsSpeed1e4(b *testing.B)    { benchmarkDecode(b, "digits-speed-1e4.br") }
func BenchmarkDecodeDigitsSpeed1e5(b *testing.B)    { benchmarkDecode(b, "digits-speed-1e5.br") }
func BenchmarkDecodeDigitsSpeed1e6(b *testing.B)    { benchmarkDecode(b, "digits-speed-1e6.br") }
func BenchmarkDecodeDigitsDefault1e4(b *testing.B)  { benchmarkDecode(b, "digits-default-1e4.br") }
func BenchmarkDecodeDigitsDefault1e5(b *testing.B)  { benchmarkDecode(b, "digits-default-1e5.br") }
func BenchmarkDecodeDigitsDefault1e6(b *testing.B)  { benchmarkDecode(b, "digits-default-1e6.br") }
func BenchmarkDecodeDigitsCompress1e4(b *testing.B) { benchmarkDecode(b, "digits-best-1e4.br") }
func BenchmarkDecodeDigitsCompress1e5(b *testing.B) { benchmarkDecode(b, "digits-best-1e5.br") }
func BenchmarkDecodeDigitsCompress1e6(b *testing.B) { benchmarkDecode(b, "digits-best-1e6.br") }
func BenchmarkDecodeTwainSpeed1e4(b *testing.B)     { benchmarkDecode(b, "twain-speed-1e4.br") }
func BenchmarkDecodeTwainSpeed1e5(b *testing.B)     { benchmarkDecode(b, "twain-speed-1e5.br") }
func BenchmarkDecodeTwainSpeed1e6(b *testing.B)     { benchmarkDecode(b, "twain-speed-1e6.br") }
func BenchmarkDecodeTwainDefault1e4(b *testing.B)   { benchmarkDecode(b, "twain-default-1e4.br") }
func BenchmarkDecodeTwainDefault1e5(b *testing.B)   { benchmarkDecode(b, "twain-default-1e5.br") }
func BenchmarkDecodeTwainDefault1e6(b *testing.B)   { benchmarkDecode(b, "twain-default-1e6.br") }
func BenchmarkDecodeTwainCompress1e4(b *testing.B)  { benchmarkDecode(b, "twain-best-1e4.br") }
func BenchmarkDecodeTwainCompress1e5(b *testing.B)  { benchmarkDecode(b, "twain-best-1e5.br") }
func BenchmarkDecodeTwainCompress1e6(b *testing.B)  { benchmarkDecode(b, "twain-best-1e6.br") }
