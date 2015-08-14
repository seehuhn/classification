package classification

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"errors"
	"github.com/seehuhn/classification/util"
	"io"
)

// TreeDecodeError is returned by `Tree.UnmarshalBinary` if the input
// is malformed.
var TreeDecodeError = errors.New("cannot decode binary tree representation")

// TreeVersionError is returned by `Tree.UnmarshalBinary` if an
// unknown encoding version is encountered.
var TreeVersionError = errors.New("unknown tree file format version")

const binaryFormatTag = "JVCT"
const binaryFormatVersion = 0

func (t *Tree) MarshalBinary() ([]byte, error) {
	buf := &bytes.Buffer{}
	err := t.WriteTo(buf)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func (t *Tree) WriteTo(w io.Writer) error {
	buf := bufio.NewWriter(w)

	// 1: tag
	_, err := buf.WriteString(binaryFormatTag)
	if err != nil {
		return err
	}

	// 2: version
	err = buf.WriteByte(binaryFormatVersion)
	if err != nil {
		return err
	}

	// 3: number of response classes
	p := t.Classes()
	err = appendUvarint(buf, uint64(p))
	if err != nil {
		return err
	}

	err = appendBinaryTree(buf, p, t)
	if err != nil {
		return err
	}

	err = buf.Flush()
	return err
}

func appendBinaryTree(buf *bufio.Writer, p int, t *Tree) error {
	if t.IsLeaf() {
		// 4: node type (0=leaf, 1=internal)
		err := buf.WriteByte(0)
		if err != nil {
			return err
		}

		// 5: histogram counts
		for i := 0; i < p; i++ {
			err = appendUvarint(buf, uint64(t.Hist[i]))
			if err != nil {
				return err
			}
		}
	} else {
		// 4: node type (0=leaf, 1=internal)
		err := buf.WriteByte(1)
		if err != nil {
			return err
		}

		// 6: split column
		err = appendUvarint(buf, uint64(t.Column))
		if err != nil {
			return err
		}

		// 7: split value
		err = binary.Write(buf, binary.LittleEndian, t.Limit)
		if err != nil {
			return err
		}

		// 8: left sub-tree
		err = appendBinaryTree(buf, p, t.LeftChild)
		if err != nil {
			return err
		}

		// 9: right sub-tree
		err = appendBinaryTree(buf, p, t.RightChild)
		if err != nil {
			return err
		}
	}
	return nil
}

func appendUvarint(buf *bufio.Writer, x uint64) error {
	tmp := [16]byte{}
	n := binary.PutUvarint(tmp[:], x)
	_, err := buf.Write(tmp[:n])
	return err
}

func (t *Tree) UnmarshalBinary(data []byte) error {
	r := bytes.NewReader(data)
	tt, err := TreeFromFile(r)
	if err != nil {
		return err
	}
	if r.Len() != 0 {
		return TreeDecodeError
	}
	*t = *tt
	return nil
}

func TreeFromFile(r io.Reader) (*Tree, error) {
	buf := bufio.NewReader(r)

	// 1: tag
	tag := make([]byte, len(binaryFormatTag))
	_, err := io.ReadFull(buf, tag)
	if err != nil {
		return nil, err
	}
	if string(tag) != binaryFormatTag {
		return nil, TreeDecodeError
	}

	// 2: version
	version, err := buf.ReadByte()
	if err != nil {
		return nil, err
	}
	if version != 0 {
		return nil, TreeVersionError
	}

	// 3: number of response classes
	pTmp, err := binary.ReadUvarint(buf)
	if err != nil {
		return nil, err
	}
	if pTmp < 1 || pTmp >= 1<<31 {
		return nil, TreeDecodeError
	}
	p := int(pTmp)

	return readBinaryTree(buf, p)
}

func readBinaryTree(buf *bufio.Reader, p int) (*Tree, error) {
	t := &Tree{}

	// 4: node type (0=leaf, 1=internal)
	nodeType, err := buf.ReadByte()
	if err != nil {
		return nil, err
	}

	if nodeType == 0 {
		// 5: histogram counts
		t.Hist = make(util.Histogram, p)
		for i := 0; i < p; i++ {
			tmp, err := binary.ReadUvarint(buf)
			if err != nil {
				return nil, err
			}
			if tmp >= 1<<31 {
				return nil, TreeDecodeError
			}
			t.Hist[i] = int(tmp)
		}
	} else if nodeType == 1 {
		// 6: split column
		tmp, err := binary.ReadUvarint(buf)
		if err != nil {
			return nil, err
		}
		if tmp >= 1<<31 {
			return nil, TreeDecodeError
		}
		t.Column = int(tmp)

		// 7: split value
		err = binary.Read(buf, binary.LittleEndian, &t.Limit)
		if err != nil {
			return nil, err
		}

		// 8: left sub-tree
		t.LeftChild, err = readBinaryTree(buf, p)
		if err != nil {
			return nil, err
		}

		// 9: right sub-tree
		t.RightChild, err = readBinaryTree(buf, p)
		if err != nil {
			return nil, err
		}

		t.Hist = make(util.Histogram, p)
		for i := 0; i < p; i++ {
			t.Hist[i] = t.LeftChild.Hist[i] + t.RightChild.Hist[i]
		}
	} else {
		return nil, TreeDecodeError
	}
	return t, nil
}
