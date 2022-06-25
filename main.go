package main

import (
    "fmt"
    "strings"
    "flag"
    "os"
    "github.com/gocolly/colly"
    "github.com/jedib0t/go-pretty/v6/table"
)

func main() {
    var ids []string
    url := "https://e3.nycu.edu.tw/"
    var session string
    flag.StringVar(&session, "session", "", "Set session cookie")
    flag.StringVar(&session, "s", "", "Set session cookie")
    var only_in_progress bool
    flag.BoolVar(&only_in_progress, "only-in-progress", false, "Only show in progress homework.")
    flag.Parse()
    if session == "" {
        fmt.Println("Session cookie must be provided.")
        os.Exit(1)
    }

    collector := colly.NewCollector()
    collector.OnResponse(func(r *colly.Response){
        //fmt.Println(string(r.Body))
    })
    collector.OnRequest(func(r *colly.Request) {
        r.Headers.Set("Cookie", "MoodleSession=" + session)
    })
    collector.OnHTML("#layer2_right_current_course_stu .course-link", func(e *colly.HTMLElement) {
        //fmt.Println(strings.Trim(e.Text, " \t"))
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
    var submitted [4][]string
    var overdue [4][]string
    if !only_in_progress{
        // submitted
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
    }
    // go visiting
    for _, id := range ids {
        collector.Visit(fmt.Sprintf(url + "local/courseextension/index.php?courseid=%s&scope=assignment", id))
    }
    // in progress
    fmt.Println("\033[33m[In Progress]\033[0m")
    t_in_progress := table.NewWriter()
    t_in_progress.SetOutputMirror(os.Stdout)
    t_in_progress.AppendHeader(table.Row{"Name", "Start", "Due", "Status"})
    for i := 0; i < len(in_progress[0]); i++ {
        name := in_progress[0][i]
        start := in_progress[1][i]
        due := in_progress[2][i]
        status := in_progress[3][i]
        t_in_progress.AppendRow([]interface{}{name, start, due, status})
    }
    t_in_progress.Render()
    if !only_in_progress {
        // submitted
        fmt.Println("\033[36m[Submitted]\033[0m")
        t_submitted := table.NewWriter()
        t_submitted.SetOutputMirror(os.Stdout)
        t_submitted.AppendHeader(table.Row{"Name", "Start", "Due", "Status"})
        for i := 0; i < len(submitted[0]); i++ {
            name := submitted[0][i]
            start := submitted[1][i]
            due := submitted[2][i]
            status := submitted[3][i]
            t_submitted.AppendRow([]interface{}{name, start, due, status})
        }
        t_submitted.Render()
        // overdue
        fmt.Println("\033[31m[Overdue]\033[0m")
        t_overdue := table.NewWriter()
        t_overdue.SetOutputMirror(os.Stdout)
        t_overdue.AppendHeader(table.Row{"Name", "Start", "Due", "Status"})
        for i := 0; i < len(overdue[0]); i++ {
            name := overdue[0][i]
            start := overdue[1][i]
            due := overdue[2][i]
            status := overdue[3][i]
            t_overdue.AppendRow([]interface{}{name, start, due, status})
        }
        t_overdue.Render()
    }
}
