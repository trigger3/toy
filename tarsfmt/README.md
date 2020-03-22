tarsfmt
--------

## 1.功能
tarsfmt能够对tars文件进行格式化，同时修复基本的错误，如结尾没有加`;`，结构体元素序标识序号错误等。

## 使用
tarsfmt的使用方法和gofmt基本相同。需要注意的是，**所有的`{`需要按照golang的语法书写，即`{`不能另起一行**。
使用效果如下

```
$ ./tarsfmt -d test
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
## 未完成功能
- 在结构体内对单独存在一行的注释的解析
- 对多行/**/ 注释的解析
- interface中目前只支持 `Resp func(ReqType req); `此种格式数据的支持
