package mytype

var Exists = struct{}{}

//set类型
type Set struct {
	// struct为结构体类型的变量
	m map[interface{}]struct{}
}

//返回一个set
func NewSet(items ...interface{}) *Set {
	// 获取Set的地址
	s := &Set{}
	// 声明map类型的数据结构
	s.m = make(map[interface{}]struct{})
	s.Add(items...)
	return s
}

func (s *Set) Size() int {
	return len(s.m)
}

//添加元素
func (s *Set) Add(items ...interface{}) {
	for _, item := range items {
		s.m[item] = Exists
	}
}

//删除元素
func (s *Set) Remove(val int) {
	delete(s.m, val)
}

//获取长度
func (s *Set) Len() int {
	return len(s.m)
}

//清空set
func (s *Set) Clear() {
	s.m = make(map[interface{}]struct{})
}

//包含
func (s *Set) Contains(item interface{}) bool {
	_, ok := s.m[item]
	return ok
}
