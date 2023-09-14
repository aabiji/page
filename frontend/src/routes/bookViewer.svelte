<script lang="ts">
    import { onMount } from "svelte";
    import { EpubViewer } from "./epub";
    import * as utils from "./utils";
    
    function getBookInfo(name: string) {
        return utils.callApi(`http://localhost:8080/book/get/${name}`, "GET", {});
    }

    let errorOut = false;
    let userBook = {
        CurrentPage: 0,
        Info: {
            Author: "",
            Contributor: "",
            Coverage: "",
            Date: "",
            Description: "",
            Identifier: "",
            Language: "",
            Publisher: "",
            Relation: "",
            Rights: "",
            Source: "",
            Subjects: [],
            Title: "",
        },
        TableOfContents: [],
        CoverImagePath: "",
    };
    onMount(() => {
        let div = document.getElementById("book-view")!;

        getBookInfo("LePetitPrince").then((json) => {
            if ("Server error" in json) {
                errorOut = true;
                console.log(json);
                return;
            }

            userBook.Info = json.Epub.Info;
            if (userBook.Info.Subjects == null)
                userBook.Info.Subjects = ["-"];

            userBook.TableOfContents = json.Epub.TableOfContents;
            userBook.CoverImagePath = utils.staticFileUrl(json.Epub.CoverImagePath);

            let e = new EpubViewer(json.FileScrollOffsets, json.Epub.Files, json.CurrentPage, div);
            e.render();
        });
    });
</script>

{#if errorOut}
<div class="error">
    <p> Oops, something went wrong </p>
</div>
{:else}
<div class="container">
    <div class="left-sidepanel">
        <h1> {userBook.Info.Title} </h1>
        <img alt="Ebook cover image" src={userBook.CoverImagePath}/>
        <h3> {userBook.Info.Author} </h3>
        <hr>
        <h5> {userBook.Info.Description} </h5>
        <p> Date: {userBook.Info.Date} </p>
        <p> Contributor: {userBook.Info.Contributor} </p>
        <p> Coverage: {userBook.Info.Coverage} </p>
        <p> Source: {userBook.Info.Source} </p>
        <p> Rights: {userBook.Info.Rights} </p>
        <p> Relation: {userBook.Info.Relation} </p>
        <p> Publisher: {userBook.Info.Publisher} </p>
        <p> Language: {userBook.Info.Language} </p>
        <p> Identifier: {userBook.Info.Identifier} </p>
        <p> Subjects: {#each userBook.Info.Subjects as subject} {subject}  {/each} </p>
        <hr>
        <h3> Table of contents </h3>
        <ol>
            {#each userBook.TableOfContents as section}
                <li><a href={section[1]}>{section[0]}</a></li>
            {/each}
        </ol>
    </div>
    <div class="right-sidepanel">
        <div id="book-view"></div>
    </div>
</div>
{/if}

<style>
    a {
        color: #4287f5;
        text-decoration: none;
    }

    .container {
        display: flex;
    }

    #book-view {
        width: 600px;
        height: 100%;
        margin: 0 auto;
    }

    .right-sidepanel {
        width: 80%;
        height: 97vh;
    }

    .left-sidepanel {
        width: 20%;
        height: 97vh;
        overflow-y: scroll;
        background-color: #1c1c1c;
        overflow-wrap: break-word;
    }

    .left-sidepanel p {
        font-size: 15px;
    }

    .left-sidepanel img {
        width: auto;
        height: 225px;
        margin-left: 25%;
    }

    .left-sidepanel h1, h3, h5 {
        text-align: center;
    }
</style>
