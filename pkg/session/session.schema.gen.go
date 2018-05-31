package session

import (
	"io"
	"time"
	"unsafe"
)

var (
	_ = unsafe.Sizeof(0)
	_ = io.ReadFull
	_ = time.Now()
)

func (d *Session) Size() (s uint64) {

	{
		l := uint64(len(d.key))

		{

			t := l
			for t >= 0x80 {
				t >>= 7
				s++
			}
			s++

		}
		s += l
	}
	s += 16
	return
}
func (d *Session) Marshal(buf []byte) ([]byte, error) {
	size := d.Size()
	{
		if uint64(cap(buf)) >= size {
			buf = buf[:size]
		} else {
			buf = make([]byte, size)
		}
	}
	i := uint64(0)

	{

		*(*uint64)(unsafe.Pointer(&buf[0])) = d.UserID

	}
	{
		l := uint64(len(d.key))

		{

			t := uint64(l)

			for t >= 0x80 {
				buf[i+8] = byte(t) | 0x80
				t >>= 7
				i++
			}
			buf[i+8] = byte(t)
			i++

		}
		copy(buf[i+8:], d.key)
		i += l
	}
	{

		*(*int64)(unsafe.Pointer(&buf[i+8])) = d.createdAt

	}
	return buf[:i+16], nil
}

func (d *Session) Unmarshal(buf []byte) (uint64, error) {
	i := uint64(0)

	{

		d.UserID = *(*uint64)(unsafe.Pointer(&buf[i+0]))

	}
	{
		l := uint64(0)

		{

			bs := uint8(7)
			t := uint64(buf[i+8] & 0x7F)
			for buf[i+8]&0x80 == 0x80 {
				i++
				t |= uint64(buf[i+8]&0x7F) << bs
				bs += 7
			}
			i++

			l = t

		}
		d.key = string(buf[i+8 : i+8+l])
		i += l
	}
	{

		d.createdAt = *(*int64)(unsafe.Pointer(&buf[i+8]))

	}
	return i + 16, nil
}
