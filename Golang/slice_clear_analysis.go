package main

import (
    "fmt"
)

func dump(letters []string) {
    fmt.Printf("addr = %p\n", letters)
    fmt.Println("letters = ", letters)
    fmt.Println(cap(letters))
    fmt.Println(len(letters))
    for i := range letters {
        fmt.Println(i, letters[i])
    }
}

func main() {
    fmt.Println("=== 基础数据 ==========")
    letters := []string{"a", "b", "c", "d"}
    dump(letters)
    
    fmt.Println("=== ====== ==========")
    
    fmt.Println("=== \"原地\"清空 ===")
    fmt.Println("=== 效果：")
    fmt.Println("=== 1.直接在原 slice 上操作，故无 GC 行为")
    fmt.Println("=== 2.清空后 cap 值和之前相同，len 值清零")
    letters = letters[:0]
    dump(letters)
    
    fmt.Println("=== 添加元素效果：基于原 slice 操作，故再未超 cap 前无需内存分配")
    letters = append(letters, "e")
    dump(letters)
    
    
    fmt.Println("=== ====== ==========")
    fmt.Println("=== 基于 nil 清空 ===")
    fmt.Println("=== 效果：")
    fmt.Println("=== 1.类似 C 语言中赋值空指针，原内容会被 GC 处理")
    fmt.Println("=== 2.清空后 cap 值清零，len 值清零")
    letters = nil
    dump(letters)
    
    fmt.Println("=== 添加元素效果：类似从无到有创建 slice")
    letters = append(letters, "e")
    dump(letters)
}
