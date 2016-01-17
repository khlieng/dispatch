package storage

// NOTE: THIS FILE WAS PRODUCED BY THE
// MSGP CODE GENERATION TOOL (github.com/tinylib/msgp)
// DO NOT EDIT

import (
	"github.com/khlieng/dispatch/Godeps/_workspace/src/github.com/tinylib/msgp/msgp"
)

// DecodeMsg implements msgp.Decodable
func (z *Channel) DecodeMsg(dc *msgp.Reader) (err error) {
	var field []byte
	_ = field
	var isz uint32
	isz, err = dc.ReadMapHeader()
	if err != nil {
		return
	}
	for isz > 0 {
		isz--
		field, err = dc.ReadMapKeyPtr()
		if err != nil {
			return
		}
		switch msgp.UnsafeString(field) {
		case "Server":
			z.Server, err = dc.ReadString()
			if err != nil {
				return
			}
		case "Name":
			z.Name, err = dc.ReadString()
			if err != nil {
				return
			}
		case "Topic":
			z.Topic, err = dc.ReadString()
			if err != nil {
				return
			}
		default:
			err = dc.Skip()
			if err != nil {
				return
			}
		}
	}
	return
}

// EncodeMsg implements msgp.Encodable
func (z Channel) EncodeMsg(en *msgp.Writer) (err error) {
	// map header, size 3
	// write "Server"
	err = en.Append(0x83, 0xa6, 0x53, 0x65, 0x72, 0x76, 0x65, 0x72)
	if err != nil {
		return err
	}
	err = en.WriteString(z.Server)
	if err != nil {
		return
	}
	// write "Name"
	err = en.Append(0xa4, 0x4e, 0x61, 0x6d, 0x65)
	if err != nil {
		return err
	}
	err = en.WriteString(z.Name)
	if err != nil {
		return
	}
	// write "Topic"
	err = en.Append(0xa5, 0x54, 0x6f, 0x70, 0x69, 0x63)
	if err != nil {
		return err
	}
	err = en.WriteString(z.Topic)
	if err != nil {
		return
	}
	return
}

// MarshalMsg implements msgp.Marshaler
func (z Channel) MarshalMsg(b []byte) (o []byte, err error) {
	o = msgp.Require(b, z.Msgsize())
	// map header, size 3
	// string "Server"
	o = append(o, 0x83, 0xa6, 0x53, 0x65, 0x72, 0x76, 0x65, 0x72)
	o = msgp.AppendString(o, z.Server)
	// string "Name"
	o = append(o, 0xa4, 0x4e, 0x61, 0x6d, 0x65)
	o = msgp.AppendString(o, z.Name)
	// string "Topic"
	o = append(o, 0xa5, 0x54, 0x6f, 0x70, 0x69, 0x63)
	o = msgp.AppendString(o, z.Topic)
	return
}

// UnmarshalMsg implements msgp.Unmarshaler
func (z *Channel) UnmarshalMsg(bts []byte) (o []byte, err error) {
	var field []byte
	_ = field
	var isz uint32
	isz, bts, err = msgp.ReadMapHeaderBytes(bts)
	if err != nil {
		return
	}
	for isz > 0 {
		isz--
		field, bts, err = msgp.ReadMapKeyZC(bts)
		if err != nil {
			return
		}
		switch msgp.UnsafeString(field) {
		case "Server":
			z.Server, bts, err = msgp.ReadStringBytes(bts)
			if err != nil {
				return
			}
		case "Name":
			z.Name, bts, err = msgp.ReadStringBytes(bts)
			if err != nil {
				return
			}
		case "Topic":
			z.Topic, bts, err = msgp.ReadStringBytes(bts)
			if err != nil {
				return
			}
		default:
			bts, err = msgp.Skip(bts)
			if err != nil {
				return
			}
		}
	}
	o = bts
	return
}

func (z Channel) Msgsize() (s int) {
	s = 1 + 7 + msgp.StringPrefixSize + len(z.Server) + 5 + msgp.StringPrefixSize + len(z.Name) + 6 + msgp.StringPrefixSize + len(z.Topic)
	return
}

// DecodeMsg implements msgp.Decodable
func (z *Message) DecodeMsg(dc *msgp.Reader) (err error) {
	var field []byte
	_ = field
	var isz uint32
	isz, err = dc.ReadMapHeader()
	if err != nil {
		return
	}
	for isz > 0 {
		isz--
		field, err = dc.ReadMapKeyPtr()
		if err != nil {
			return
		}
		switch msgp.UnsafeString(field) {
		case "ID":
			z.ID, err = dc.ReadUint64()
			if err != nil {
				return
			}
		case "Server":
			z.Server, err = dc.ReadString()
			if err != nil {
				return
			}
		case "From":
			z.From, err = dc.ReadString()
			if err != nil {
				return
			}
		case "To":
			z.To, err = dc.ReadString()
			if err != nil {
				return
			}
		case "Content":
			z.Content, err = dc.ReadString()
			if err != nil {
				return
			}
		case "Time":
			z.Time, err = dc.ReadInt64()
			if err != nil {
				return
			}
		default:
			err = dc.Skip()
			if err != nil {
				return
			}
		}
	}
	return
}

// EncodeMsg implements msgp.Encodable
func (z *Message) EncodeMsg(en *msgp.Writer) (err error) {
	// map header, size 6
	// write "ID"
	err = en.Append(0x86, 0xa2, 0x49, 0x44)
	if err != nil {
		return err
	}
	err = en.WriteUint64(z.ID)
	if err != nil {
		return
	}
	// write "Server"
	err = en.Append(0xa6, 0x53, 0x65, 0x72, 0x76, 0x65, 0x72)
	if err != nil {
		return err
	}
	err = en.WriteString(z.Server)
	if err != nil {
		return
	}
	// write "From"
	err = en.Append(0xa4, 0x46, 0x72, 0x6f, 0x6d)
	if err != nil {
		return err
	}
	err = en.WriteString(z.From)
	if err != nil {
		return
	}
	// write "To"
	err = en.Append(0xa2, 0x54, 0x6f)
	if err != nil {
		return err
	}
	err = en.WriteString(z.To)
	if err != nil {
		return
	}
	// write "Content"
	err = en.Append(0xa7, 0x43, 0x6f, 0x6e, 0x74, 0x65, 0x6e, 0x74)
	if err != nil {
		return err
	}
	err = en.WriteString(z.Content)
	if err != nil {
		return
	}
	// write "Time"
	err = en.Append(0xa4, 0x54, 0x69, 0x6d, 0x65)
	if err != nil {
		return err
	}
	err = en.WriteInt64(z.Time)
	if err != nil {
		return
	}
	return
}

// MarshalMsg implements msgp.Marshaler
func (z *Message) MarshalMsg(b []byte) (o []byte, err error) {
	o = msgp.Require(b, z.Msgsize())
	// map header, size 6
	// string "ID"
	o = append(o, 0x86, 0xa2, 0x49, 0x44)
	o = msgp.AppendUint64(o, z.ID)
	// string "Server"
	o = append(o, 0xa6, 0x53, 0x65, 0x72, 0x76, 0x65, 0x72)
	o = msgp.AppendString(o, z.Server)
	// string "From"
	o = append(o, 0xa4, 0x46, 0x72, 0x6f, 0x6d)
	o = msgp.AppendString(o, z.From)
	// string "To"
	o = append(o, 0xa2, 0x54, 0x6f)
	o = msgp.AppendString(o, z.To)
	// string "Content"
	o = append(o, 0xa7, 0x43, 0x6f, 0x6e, 0x74, 0x65, 0x6e, 0x74)
	o = msgp.AppendString(o, z.Content)
	// string "Time"
	o = append(o, 0xa4, 0x54, 0x69, 0x6d, 0x65)
	o = msgp.AppendInt64(o, z.Time)
	return
}

// UnmarshalMsg implements msgp.Unmarshaler
func (z *Message) UnmarshalMsg(bts []byte) (o []byte, err error) {
	var field []byte
	_ = field
	var isz uint32
	isz, bts, err = msgp.ReadMapHeaderBytes(bts)
	if err != nil {
		return
	}
	for isz > 0 {
		isz--
		field, bts, err = msgp.ReadMapKeyZC(bts)
		if err != nil {
			return
		}
		switch msgp.UnsafeString(field) {
		case "ID":
			z.ID, bts, err = msgp.ReadUint64Bytes(bts)
			if err != nil {
				return
			}
		case "Server":
			z.Server, bts, err = msgp.ReadStringBytes(bts)
			if err != nil {
				return
			}
		case "From":
			z.From, bts, err = msgp.ReadStringBytes(bts)
			if err != nil {
				return
			}
		case "To":
			z.To, bts, err = msgp.ReadStringBytes(bts)
			if err != nil {
				return
			}
		case "Content":
			z.Content, bts, err = msgp.ReadStringBytes(bts)
			if err != nil {
				return
			}
		case "Time":
			z.Time, bts, err = msgp.ReadInt64Bytes(bts)
			if err != nil {
				return
			}
		default:
			bts, err = msgp.Skip(bts)
			if err != nil {
				return
			}
		}
	}
	o = bts
	return
}

func (z *Message) Msgsize() (s int) {
	s = 1 + 3 + msgp.Uint64Size + 7 + msgp.StringPrefixSize + len(z.Server) + 5 + msgp.StringPrefixSize + len(z.From) + 3 + msgp.StringPrefixSize + len(z.To) + 8 + msgp.StringPrefixSize + len(z.Content) + 5 + msgp.Int64Size
	return
}

// DecodeMsg implements msgp.Decodable
func (z *Server) DecodeMsg(dc *msgp.Reader) (err error) {
	var field []byte
	_ = field
	var isz uint32
	isz, err = dc.ReadMapHeader()
	if err != nil {
		return
	}
	for isz > 0 {
		isz--
		field, err = dc.ReadMapKeyPtr()
		if err != nil {
			return
		}
		switch msgp.UnsafeString(field) {
		case "Name":
			z.Name, err = dc.ReadString()
			if err != nil {
				return
			}
		case "Host":
			z.Host, err = dc.ReadString()
			if err != nil {
				return
			}
		case "Port":
			z.Port, err = dc.ReadString()
			if err != nil {
				return
			}
		case "TLS":
			z.TLS, err = dc.ReadBool()
			if err != nil {
				return
			}
		case "Password":
			z.Password, err = dc.ReadString()
			if err != nil {
				return
			}
		case "Nick":
			z.Nick, err = dc.ReadString()
			if err != nil {
				return
			}
		case "Username":
			z.Username, err = dc.ReadString()
			if err != nil {
				return
			}
		case "Realname":
			z.Realname, err = dc.ReadString()
			if err != nil {
				return
			}
		default:
			err = dc.Skip()
			if err != nil {
				return
			}
		}
	}
	return
}

// EncodeMsg implements msgp.Encodable
func (z *Server) EncodeMsg(en *msgp.Writer) (err error) {
	// map header, size 8
	// write "Name"
	err = en.Append(0x88, 0xa4, 0x4e, 0x61, 0x6d, 0x65)
	if err != nil {
		return err
	}
	err = en.WriteString(z.Name)
	if err != nil {
		return
	}
	// write "Host"
	err = en.Append(0xa4, 0x48, 0x6f, 0x73, 0x74)
	if err != nil {
		return err
	}
	err = en.WriteString(z.Host)
	if err != nil {
		return
	}
	// write "Port"
	err = en.Append(0xa4, 0x50, 0x6f, 0x72, 0x74)
	if err != nil {
		return err
	}
	err = en.WriteString(z.Port)
	if err != nil {
		return
	}
	// write "TLS"
	err = en.Append(0xa3, 0x54, 0x4c, 0x53)
	if err != nil {
		return err
	}
	err = en.WriteBool(z.TLS)
	if err != nil {
		return
	}
	// write "Password"
	err = en.Append(0xa8, 0x50, 0x61, 0x73, 0x73, 0x77, 0x6f, 0x72, 0x64)
	if err != nil {
		return err
	}
	err = en.WriteString(z.Password)
	if err != nil {
		return
	}
	// write "Nick"
	err = en.Append(0xa4, 0x4e, 0x69, 0x63, 0x6b)
	if err != nil {
		return err
	}
	err = en.WriteString(z.Nick)
	if err != nil {
		return
	}
	// write "Username"
	err = en.Append(0xa8, 0x55, 0x73, 0x65, 0x72, 0x6e, 0x61, 0x6d, 0x65)
	if err != nil {
		return err
	}
	err = en.WriteString(z.Username)
	if err != nil {
		return
	}
	// write "Realname"
	err = en.Append(0xa8, 0x52, 0x65, 0x61, 0x6c, 0x6e, 0x61, 0x6d, 0x65)
	if err != nil {
		return err
	}
	err = en.WriteString(z.Realname)
	if err != nil {
		return
	}
	return
}

// MarshalMsg implements msgp.Marshaler
func (z *Server) MarshalMsg(b []byte) (o []byte, err error) {
	o = msgp.Require(b, z.Msgsize())
	// map header, size 8
	// string "Name"
	o = append(o, 0x88, 0xa4, 0x4e, 0x61, 0x6d, 0x65)
	o = msgp.AppendString(o, z.Name)
	// string "Host"
	o = append(o, 0xa4, 0x48, 0x6f, 0x73, 0x74)
	o = msgp.AppendString(o, z.Host)
	// string "Port"
	o = append(o, 0xa4, 0x50, 0x6f, 0x72, 0x74)
	o = msgp.AppendString(o, z.Port)
	// string "TLS"
	o = append(o, 0xa3, 0x54, 0x4c, 0x53)
	o = msgp.AppendBool(o, z.TLS)
	// string "Password"
	o = append(o, 0xa8, 0x50, 0x61, 0x73, 0x73, 0x77, 0x6f, 0x72, 0x64)
	o = msgp.AppendString(o, z.Password)
	// string "Nick"
	o = append(o, 0xa4, 0x4e, 0x69, 0x63, 0x6b)
	o = msgp.AppendString(o, z.Nick)
	// string "Username"
	o = append(o, 0xa8, 0x55, 0x73, 0x65, 0x72, 0x6e, 0x61, 0x6d, 0x65)
	o = msgp.AppendString(o, z.Username)
	// string "Realname"
	o = append(o, 0xa8, 0x52, 0x65, 0x61, 0x6c, 0x6e, 0x61, 0x6d, 0x65)
	o = msgp.AppendString(o, z.Realname)
	return
}

// UnmarshalMsg implements msgp.Unmarshaler
func (z *Server) UnmarshalMsg(bts []byte) (o []byte, err error) {
	var field []byte
	_ = field
	var isz uint32
	isz, bts, err = msgp.ReadMapHeaderBytes(bts)
	if err != nil {
		return
	}
	for isz > 0 {
		isz--
		field, bts, err = msgp.ReadMapKeyZC(bts)
		if err != nil {
			return
		}
		switch msgp.UnsafeString(field) {
		case "Name":
			z.Name, bts, err = msgp.ReadStringBytes(bts)
			if err != nil {
				return
			}
		case "Host":
			z.Host, bts, err = msgp.ReadStringBytes(bts)
			if err != nil {
				return
			}
		case "Port":
			z.Port, bts, err = msgp.ReadStringBytes(bts)
			if err != nil {
				return
			}
		case "TLS":
			z.TLS, bts, err = msgp.ReadBoolBytes(bts)
			if err != nil {
				return
			}
		case "Password":
			z.Password, bts, err = msgp.ReadStringBytes(bts)
			if err != nil {
				return
			}
		case "Nick":
			z.Nick, bts, err = msgp.ReadStringBytes(bts)
			if err != nil {
				return
			}
		case "Username":
			z.Username, bts, err = msgp.ReadStringBytes(bts)
			if err != nil {
				return
			}
		case "Realname":
			z.Realname, bts, err = msgp.ReadStringBytes(bts)
			if err != nil {
				return
			}
		default:
			bts, err = msgp.Skip(bts)
			if err != nil {
				return
			}
		}
	}
	o = bts
	return
}

func (z *Server) Msgsize() (s int) {
	s = 1 + 5 + msgp.StringPrefixSize + len(z.Name) + 5 + msgp.StringPrefixSize + len(z.Host) + 5 + msgp.StringPrefixSize + len(z.Port) + 4 + msgp.BoolSize + 9 + msgp.StringPrefixSize + len(z.Password) + 5 + msgp.StringPrefixSize + len(z.Nick) + 9 + msgp.StringPrefixSize + len(z.Username) + 9 + msgp.StringPrefixSize + len(z.Realname)
	return
}

// DecodeMsg implements msgp.Decodable
func (z *User) DecodeMsg(dc *msgp.Reader) (err error) {
	var field []byte
	_ = field
	var isz uint32
	isz, err = dc.ReadMapHeader()
	if err != nil {
		return
	}
	for isz > 0 {
		isz--
		field, err = dc.ReadMapKeyPtr()
		if err != nil {
			return
		}
		switch msgp.UnsafeString(field) {
		case "ID":
			z.ID, err = dc.ReadUint64()
			if err != nil {
				return
			}
		case "Username":
			z.Username, err = dc.ReadString()
			if err != nil {
				return
			}
		default:
			err = dc.Skip()
			if err != nil {
				return
			}
		}
	}
	return
}

// EncodeMsg implements msgp.Encodable
func (z *User) EncodeMsg(en *msgp.Writer) (err error) {
	// map header, size 2
	// write "ID"
	err = en.Append(0x82, 0xa2, 0x49, 0x44)
	if err != nil {
		return err
	}
	err = en.WriteUint64(z.ID)
	if err != nil {
		return
	}
	// write "Username"
	err = en.Append(0xa8, 0x55, 0x73, 0x65, 0x72, 0x6e, 0x61, 0x6d, 0x65)
	if err != nil {
		return err
	}
	err = en.WriteString(z.Username)
	if err != nil {
		return
	}
	return
}

// MarshalMsg implements msgp.Marshaler
func (z *User) MarshalMsg(b []byte) (o []byte, err error) {
	o = msgp.Require(b, z.Msgsize())
	// map header, size 2
	// string "ID"
	o = append(o, 0x82, 0xa2, 0x49, 0x44)
	o = msgp.AppendUint64(o, z.ID)
	// string "Username"
	o = append(o, 0xa8, 0x55, 0x73, 0x65, 0x72, 0x6e, 0x61, 0x6d, 0x65)
	o = msgp.AppendString(o, z.Username)
	return
}

// UnmarshalMsg implements msgp.Unmarshaler
func (z *User) UnmarshalMsg(bts []byte) (o []byte, err error) {
	var field []byte
	_ = field
	var isz uint32
	isz, bts, err = msgp.ReadMapHeaderBytes(bts)
	if err != nil {
		return
	}
	for isz > 0 {
		isz--
		field, bts, err = msgp.ReadMapKeyZC(bts)
		if err != nil {
			return
		}
		switch msgp.UnsafeString(field) {
		case "ID":
			z.ID, bts, err = msgp.ReadUint64Bytes(bts)
			if err != nil {
				return
			}
		case "Username":
			z.Username, bts, err = msgp.ReadStringBytes(bts)
			if err != nil {
				return
			}
		default:
			bts, err = msgp.Skip(bts)
			if err != nil {
				return
			}
		}
	}
	o = bts
	return
}

func (z *User) Msgsize() (s int) {
	s = 1 + 3 + msgp.Uint64Size + 9 + msgp.StringPrefixSize + len(z.Username)
	return
}
