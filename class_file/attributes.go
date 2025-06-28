package class_file

import "github.com/takeru-a/self-jvm/support"

type (
	// Code属性を表す構造体
	CodeAttr struct {
		maxStack  uint16 // オペランドスタックのサイズ
		maxLocals uint16 // ローカル変数領域のサイズ
		code      []byte // 命令
	}
)

const (
	// Code属性を示す文字列
	codeAttr = "Code"
)

// 各種attributesフィールドをデコード。
// 属性の種別を判断するため定数プールにアクセスする必要があるので、引数で定数プールを表す構造体ConstantPoolを取る。
func readAttributes(r *support.ByteSeq, cp *ConstantPool) []interface{} {
	// attributes_conuntフィールドのデコード
	attrs := make([]interface{}, r.ReadUint16())
	for i := 0; i < len(attrs); i++ {
		attrs[i] = readAttribute(r, cp)
	}

	return attrs
}

func readAttribute(r *support.ByteSeq, cp *ConstantPool) interface{} {
	// attribute_name_indexフィールドから定数プール中のインデックスをデコードし、既に実装したUtf8メソッドで属性の名前を取得。
	// u2 attribute_name_index;
	name := cp.Utf8(r.ReadUint16())
	// attribute_lengthフィールドのデコード
	// u4 attribute_length
	size := r.ReadUint32()

	switch name {
	case codeAttr:
		return readCodeAttr(r, cp)
	default:
		// Code属性以外の属性については、
		// 詳細を知らなくてもサイズが分かっているため単にスキップ可能
		r.Skip(int(size))
		return nil
	}
}

// Code属性のデコード
func readCodeAttr(r *support.ByteSeq, cp *ConstantPool) interface{} {
	// max_stack、max_locals、codeフィールドをデコード
	attr := &CodeAttr{
		maxStack:  r.ReadUint16(),
		maxLocals: r.ReadUint16(),
		code:      r.ReadBytes(int(r.ReadUint32())),
	}
	// exception_tableフィールドをスキップ。
	// (exception_table_lengthフィールドが示す分、u2フィールド4つ=8バイトをスキップ)
	r.Skip(int(r.ReadUint16() * 8))
	// 属性をスキップ(デコードするが戻り値を捨てる)
	readAttributes(r, cp)
	return attr
}

// Code属性が持つローカル変数領域のサイズを返す。
func (ca *CodeAttr) MaxLocals() uint16 {
	return ca.maxLocals
}

// Code属性が持つ命令を返す。
func (ca *CodeAttr) Code() []byte {
	return ca.code
}
