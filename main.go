package main

import (
    "fmt"
    "strings"
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
    })

    collector.Visit(url)

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

    // submitted
    var submitted [4][]string
    collector.OnHTML("#news-view-nofile2-tobegraded-in-progress tbody tr .instancename", func(e *colly.HTMLElement) {
        // HW name
        submitted[0] = append(submitted[0], strings.Trim(e.Text, " \t"))
    })
    collector.OnHTML("#news-view-nofile2-tobegraded-in-progress tbody tr td:nth-child(2)", func(e *colly.HTMLElement) {
        // start date
        submitted[1] = append(submitted[1], strings.Trim(e.Text, " \t"))
    })
    collector.OnHTML("#news-view-nofile2-tobegraded-in-progress tbody tr td:nth-child(3)", func(e *colly.HTMLElement) {
        // due date
        submitted[2] = append(submitted[2], strings.Trim(e.Text, " \t"))
    })
    collector.OnHTML("#news-view-nofile2-tobegraded-in-progress tbody tr td:nth-child(4)", func(e *colly.HTMLElement) {
        // people count
        submitted[3] = append(submitted[3], strings.Trim(e.Text, " \t"))
    })
    
    // overdue
    var overdue [4][]string
    collector.OnHTML("#news-view-nofile2-notsubmitted-in-progress tbody tr .instancename", func(e *colly.HTMLElement) {
        // HW name
        overdue[0] = append(overdue[0], strings.Trim(e.Text, " \t"))
    })
    collector.OnHTML("#news-view-nofile2-notsubmitted-in-progress tbody tr td:nth-child(2)", func(e *colly.HTMLElement) {
        // start date
        overdue[1] = append(overdue[1], strings.Trim(e.Text, " \t"))
    })
    collector.OnHTML("#news-view-nofile2-notsubmitted-in-progress tbody tr td:nth-child(3)", func(e *colly.HTMLElement) {
        // due date
        overdue[2] = append(overdue[2], strings.Trim(e.Text, " \t"))
    })
    collector.OnHTML("#news-view-nofile2-notsubmitted-in-progress tbody tr td:nth-child(4)", func(e *colly.HTMLElement) {
        // people count
        overdue[3] = append(overdue[3], strings.Trim(e.Text, " \t"))
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
    fmt.Println("---------------\n")
    // submitted
    fmt.Println("[Submitted]")
    for i := 0; i < len(submitted[0]); i++ {
        fmt.Println("Name: ", submitted[0][i])
        fmt.Println("Start: ", submitted[1][i])
        fmt.Println("Due: ", submitted[2][i])
        fmt.Println("Status: ", submitted[3][i], "\n")
    }
    fmt.Println("---------------\n")
    // overdue
    fmt.Println("[Overdue]")
    for i := 0; i < len(overdue[0]); i++ {
        fmt.Println("Name: ", overdue[0][i])
        fmt.Println("Start: ", overdue[1][i])
        fmt.Println("Due: ", overdue[2][i])
        fmt.Println("Status: ", overdue[3][i], "\n")
    }
}
