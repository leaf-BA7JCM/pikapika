package main

import (
        "ci/commons"
        "encoding/json"
        "fmt"
        "io/ioutil"
        "net/http"
        "os"
        "strings"
)

const owner = "niuhuan"
const repo = "pikapika"
const ua = "niuhuan pikapika ci"

func main() {
        // get ghToken
        ghToken := os.Getenv("GH_TOKEN")
        if ghToken == "" {
                println("Env ${GH_TOKEN} is not set")
                os.Exit(1)
        }
        // get version
        var version commons.Version
        codeFile, err := ioutil.ReadFile("version.code.txt")
        if err != nil {
                panic(err)
        }
        version.Code = strings.TrimSpace(string(codeFile))
        infoFile, err := ioutil.ReadFile("version.info.txt")
        if err != nil {
                panic(err)
        }
        version.Info = strings.TrimSpace(string(infoFile))
        // get target
        target := os.Getenv("TARGET")
        if ghToken == "" {
                println("Env ${TARGET} is not set")
                os.Exit(1)
        }
        //
        var releaseFileName string
        switch target {
        case "macos":
                releaseFileName = fmt.Sprintf("pikapika-v1.4.1-android-arm32.apk")
        case "ios":
                releaseFileName = fmt.Sprintf("pikapika-%v-ios-nosign.ipa", version.Code)
        case "windows":
                releaseFileName = fmt.Sprintf("pikapika-%v-windows-x86_64.zip", version.Code)
        case "linux":
                releaseFileName = fmt.Sprintf("pikapika-%v-linux-x86_64.AppImage", version.Code)
        case "android-arm32":
                releaseFileName = fmt.Sprintf("pikapika-%v-android-arm32.apk", version.Code)
        case "android-arm64":
                releaseFileName = fmt.Sprintf("pikapika-%v-android-arm64.apk", version.Code)
        case "android-x86_64":
                releaseFileName = fmt.Sprintf("pikapika-%v-android-x86_64.apk", version.Code)
        }
        // get version
        getReleaseRequest, err := http.NewRequest(
                "GET",
                fmt.Sprintf("https://api.github.com/repos/%v/%v/releases/tags/%v", owner, repo, version.Code),
                nil,
        )
        if err != nil {
                panic(err)
        }
        getReleaseRequest.Header.Set("User-Agent", ua)
        getReleaseRequest.Header.Set("Authorization", ghToken)
        getReleaseResponse, err := http.DefaultClient.Do(getReleaseRequest)
        if err != nil {
                panic(err)
        }
        defer getReleaseResponse.Body.Close()
        if getReleaseResponse.StatusCode == 404 {
                panic("NOT FOUND RELEASE")
        }
        buff, err := ioutil.ReadAll(getReleaseResponse.Body)
        if err != nil {
                panic(err)
        }
        var release commons.Release
        err = json.Unmarshal(buff, &release)
        if err != nil {
                println(string(buff))
                panic(err)
        }
        for _, asset := range release.Assets {
                if asset.Name == releaseFileName {
                        println("::set-output name=skip_build::true")
                        os.Exit(0)
                }
        }
        print("::set-output name=skip_build::false")
}
