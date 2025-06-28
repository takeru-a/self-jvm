package class_file

import (
	"github.com/takeru-a/self-jvm/support"
)

type (
	// 定数プールの定義
	ConstantPool struct {
		cpInfo []interface{} // 文字列、クラス、メソッド、シグネチャなどの定数プール情報
	}

	// CONSTANT_Class
	ClassCpInfo uint16 // クラス名

	// CONSTANT_Methodref メソッドの定義
	MethodRefCpInfo struct {
		class       uint16 // クラス名
		nameAndType uint16 // メソッド名とシグネチャ
	}

	// CONSTANT_NameAndType
	NameAndTypeCpInfo struct {
		name uint16 // メソッド名
		desc uint16 // メソッドのシグネチャ
	}
)

// 定数プールのタグ値
const (
	utf8Tag        uint8 = 1 // 文字列定数
	classTag       uint8 = 7 // クラス
	methodRefTag   uint8 = 10 // メソッド
	nameAndTypeTag uint8 = 12 // メソッド名とシグネチャ
)

// constant_pool_countフィールドの直前まで読み進めた構造体ByteSeqを受け取り、定数プールをデコード。
func readCP(r *support.ByteSeq) *ConstantPool {
	cpCount := r.ReadUint16()
	// 仕様上、constant_poolフィールドの長さはconstant_pool_count-1となっているため
	// 1減じる必要があるが、インデックス番号が1から始まることもあり、
	// スライスの長さはconstant_pool_countとし、スライスの2番目の要素から定数プールエントリーを格納していく。
	cp := &ConstantPool{cpInfo: make([]interface{}, cpCount)}
	for i := uint16(1); i < cpCount; i++ {
		switch r.ReadByte() {
		case utf8Tag:
			// CONSTANT_Utf8用の構造体は定義せず、単にstringとする
			cp.cpInfo[i] = string(r.ReadBytes(int(r.ReadUint16())))
		case classTag:
			cp.cpInfo[i] = ClassCpInfo(r.ReadUint16())
		case methodRefTag:
			cp.cpInfo[i] = &MethodRefCpInfo{class: r.ReadUint16(), nameAndType: r.ReadUint16()}
		case nameAndTypeTag:
			cp.cpInfo[i] = &NameAndTypeCpInfo{name: r.ReadUint16(),
				desc: r.ReadUint16(),
			}
		}
	}
	return cp
}

//  定数プール中のCONSTANT_Utf8エントリーを返す
func (cp *ConstantPool) Utf8(index uint16) string {
	s, ok := cp.cpInfo[index].(string)
	if !ok {
		return ""
	}
	return s
}
