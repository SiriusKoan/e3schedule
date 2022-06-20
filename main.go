package main

import (
    "fmt"
    //"os"
    "flag"
    "github.com/gocolly/colly"
)

func main() {
    var update bool
    flag.BoolVar(&update, "update", false, "force update or not")
    flag.BoolVar(&update, "u", false, "force update or not")
    var session string
    flag.StringVar(&session, "session", "", "session cookie")
    flag.StringVar(&session, "s", "", "session cookie")
    flag.Parse()

    collector := colly.NewCollector()
    collector.OnResponse(func(r *colly.Response){
        //fmt.Println(string(r.Body))
    })
    collector.OnRequest(func(r *colly.Request) {
        r.Headers.Set("Cookie", "MoodleSession=" + session)
    })
    collector.OnHTML("#layer2_right_current_course_stu .course-link", func(e *colly.HTMLElement) {
        fmt.Println(e.Attr("href"))
    })

    collector.Visit("https://e3.nycu.edu.tw")

}
