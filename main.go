package main

import (
    "fmt"
    "strings"
    //"os"
    "flag"
    "github.com/gocolly/colly"
)

func main() {
    var ids []string
    url := "https://e3.nycu.edu.tw/"
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
        fmt.Println(strings.Trim(e.Text, " \t"))
        ids = append(ids, strings.Split(e.Attr("href"), "id=")[1])
        //collector.Visi)(fmt.Sprintf(url + "local/courseextension/index.php?courseid=%s&scope=assignment", id))
    })
    /*
    */

    collector.Visit(url)

    //cur := 0;
    // in progress
    var in_progress [4][]string
    collector.OnHTML("#news-view-basic-in-progress tbody tr .instancename", func(e *colly.HTMLElement) {
        // HW name
        in_progress[0] = append(in_progress[0], strings.Trim(e.Text, " \t"))
    })
    collector.OnHTML("#news-view-basic-in-progress tbody tr td:nth-child(2)", func(e *colly.HTMLElement) {
        // start date
        in_progress[1] = append(in_progress[1], strings.Trim(e.Text, " \t"))
    })
    collector.OnHTML("#news-view-basic-in-progress tbody tr td:nth-child(3)", func(e *colly.HTMLElement) {
        // due date
        in_progress[2] = append(in_progress[2], strings.Trim(e.Text, " \t"))
    })
    collector.OnHTML("#news-view-basic-in-progress tbody tr td:nth-child(4)", func(e *colly.HTMLElement) {
        // people count
        in_progress[3] = append(in_progress[3], strings.Trim(e.Text, " \t"))
    })

    // done
    var done [4][]string
    collector.OnHTML("#news-view-nofile2-tobegraded-in-progress tbody tr .instancename", func(e *colly.HTMLElement) {
        // HW name
        done[0] = append(done[0], strings.Trim(e.Text, " \t"))
    })
    collector.OnHTML("#news-view-nofile2-tobegraded-in-progress tbody tr td:nth-child(2)", func(e *colly.HTMLElement) {
        // start date
        done[1] = append(done[1], strings.Trim(e.Text, " \t"))
    })
    collector.OnHTML("#news-view-nofile2-tobegraded-in-progress tbody tr td:nth-child(3)", func(e *colly.HTMLElement) {
        // due date
        done[2] = append(done[2], strings.Trim(e.Text, " \t"))
    })
    collector.OnHTML("#news-view-nofile2-tobegraded-in-progress tbody tr td:nth-child(4)", func(e *colly.HTMLElement) {
        // people count
        done[3] = append(done[3], strings.Trim(e.Text, " \t"))
    })
    /*
    collector.OnHTML("#news-view-nofile2-notsubmitted-in-progress tbody tr", func(e *colly.HTMLElement) {
        fmt.Println(e.Text)
    })
    */
    // go visiting
    for _, id := range ids {
        collector.Visit(fmt.Sprintf(url + "local/courseextension/index.php?courseid=%s&scope=assignment", id))
    }
    // in progress
    fmt.Println("[In Progress]")
    for i := 0; i < len(in_progress[0]); i++ {
        fmt.Println("Name: ", in_progress[0][i])
        fmt.Println("Start: ", in_progress[1][i])
        fmt.Println("Due: ", in_progress[2][i])
        fmt.Println("Status: ", in_progress[3][i], "\n")
    }
    fmt.Println("---------------")
    // done
    fmt.Println("[Done (Submitted)]")
    for i := 0; i < len(done[0]); i++ {
        fmt.Println("Name: ", done[0][i])
        fmt.Println("Start: ", done[1][i])
        fmt.Println("Due: ", done[2][i])
        fmt.Println("Status: ", done[3][i], "\n")
    }
}
