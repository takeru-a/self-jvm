package class_file

import (
	"fmt"
	"io"

	"github.com/takeru-a/self-jvm/support"
)

type (
	ClassFile struct {
		cp      *ConstantPool // 定数プール
		methods []*MethodInfo //メソッド
	}

	// methodの情報を保持する構造体
	MethodInfo struct {
		accessFlag uint16
		name       string
		desc       string
		attributes []interface{}
	}
)

// JVM仕様で定義されるclassファイルの固定マジックナンバー
const magicNumber = 0xCAFEBABE

func ReadClassFile(cfReader io.Reader) (*ClassFile, error) {
	r, err := support.NewByteSeq(cfReader)
	if err != nil {
		return nil, err
	}
	//  マジックナンバー(u4  =  非負4バイト整数)を検証
	// u4 magic
	if r.ReadUint32() != magicNumber {
		return nil, fmt.Errorf("not  java  class  file")
	}
	// マイナー・メジャーバージョン(u2  =  非負2バイト整数)は単にスキップ(合計4バイト)
	// u2 minor_version
	// u2 major_version
	r.Skip(4)

	// 定数プールを読み取り
	// cp_info constant_pool_count
	class := &ClassFile{cp: readCP(r)}

	// access_flags,  this_class,  super_classフィールドのスキップ
	// u2 access_flags pubulic、private、final等
	// u2 this_class   自身のクラス
	// u2 super_class  親のクラス
	r.Skip(6)

	// interfaces_conutフィールドをデコードし、その数だけinterfacesフィールドのデコードをスキップする
	// u2 interfaces_count
	// u2 interfaces[ineterfaces_count]
	r.Skip(int(r.ReadUint16() * 2))

	// クラスのフィールドを扱わない。
	// クラスのフィールド(fieldsフィールド)はメソッド(methodsフィールド)と同じ構造なので、同様のメソッドを呼び出してデコードをスキップする。
	readMethodInfo(r, class.cp)

	class.methods = readMethodInfo(r, class.cp)

	// methodsフィールドの後には本来属性関連のフィールドが続くが、method_infoの属性であるCode属性以外を扱わないので、デコードせず処理を返す。

	return class, nil
}

// methodsフィールドのデコード
func readMethodInfo(r *support.ByteSeq, cp *ConstantPool) []*MethodInfo {
	count := r.ReadUint16()
	methods := make([]*MethodInfo, count)

	for i := uint16(0); i < count; i++ {
		methods[i] = &MethodInfo{
			accessFlag: r.ReadUint16(),
			name:       cp.Utf8(r.ReadUint16()),
			desc:       cp.Utf8(r.ReadUint16()),
			attributes: readAttributes(r, cp),
		}
	}

	return methods
}

// 保持する属性のうちCode属性を返す
func (m *MethodInfo) Code() *CodeAttr {
	for _, attr := range m.attributes {
		if code, ok := attr.(*CodeAttr); ok {
			return code
		}
	}
	return nil
}

// デコードしたメソッド一覧を返す
func (c *ClassFile) Methods() []*MethodInfo {
	return c.methods
}

// メソッドの文字列表現を返す(名前とシグネチャ)
func (m *MethodInfo) String() string {
	return m.name + m.desc
}
