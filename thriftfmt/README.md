thriftfmt
--------

## 1.功能
thriftfmt能够对tars文件进行格式化，同时修复基本的错误，如结尾没有加`;`，结构体元素序标识序号错误等。

## 2.使用
1. thriftfmt的使用方法和gofmt基本相同。
2. 注意如下两点
- 所有的`{`需要按照golang的语法书写，即`{`不能另起一行；
- 自定义元素在引用时，其定义语句在文件必须在引用点的上部；

效果如下：
```
$ ./thriftfmt -d test
diff -u test/test.tars.orig test/test.tars
--- test/test.tars.orig	2020-03-22 12:41:42.000000000 +0800
+++ test/test.tars	2020-03-22 12:41:42.000000000 +0800
@@ -4,14 +4,15 @@
     };

     struct Value {
-        0 require int a; // 测试对注释的格式化
-        0 require unsigned int b;// 测试对元素标识符的修正
-        0 require unsigned int c//test 测试对确实结束符”;“的修正
-        0 require unsigned int d;//test
+        0 require int a;          // 测试对注释的格式化
+        1 require unsigned int b; // 测试对元素标识符的修正
+        2 require unsigned int c; // test 测试对确实结束符”;“的修正
+        3 require unsigned int d; // test
     };

     interface TestLogic { // test
-        int Value (Key k);//testsda
-        int Value (Key k321321);//t3wst
+        int Value(Key k);       // testsda
+        int Value(Key k321321); // t3wst
     };
 };
+
```
## 3. TODOLIST
- 对多行/**/ 注释的解析；
- interface中目前只支持 `Resp func(ReqType req); `此种格式数据的支持；
