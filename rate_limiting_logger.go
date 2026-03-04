package main

type Logger struct {
    lastSeen map[string]int
}

func Constructor() Logger {
    return Logger {
        lastSeen: make(map[string]int),
    }
}

func (l *Logger) ShouldPrintMessage(timestamp int, message string) bool {
    if lastTime, exists := l.lastSeen[message]; exists {
        if timestamp - lastTime < 10 {
            return false
        }
    }
    l.lastSeen[message] = timestamp
    return true
}


func main() {
    Logger := Constructor()
    println(Logger.ShouldPrintMessage(1, "foo")) // true
    println(Logger.ShouldPrintMessage(2, "bar")) // true
    println(Logger.ShouldPrintMessage(3, "foo")) // false
    println(Logger.ShouldPrintMessage(8, "bar")) // false
    println(Logger.ShouldPrintMessage(10, "foo")) // false
    println(Logger.ShouldPrintMessage(11, "foo")) // true
}