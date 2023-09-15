<script lang="ts">
    import { onMount } from "svelte";
    import { writable } from "svelte/store";

    import { EpubViewer } from "./epub";
    import * as utils from "./utils";

    let errorOut = false;
    let book = writable({
        CurrentPage: 0,
        Epub: {
            CoverImagePath: "",
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
        },
    });

    function getBook(name: string, div: HTMLElement) {
        utils.callApi(`http://localhost:8080/book/get/${name}`, "GET", {}).then((json) => {
            if ("Server error" in json) {
                errorOut = true;
                console.log(json);
                return;
            }

            json.Epub.CoverImagePath = utils.staticFileUrl(json.Epub.CoverImagePath);
            book.set(json);

            let e = new EpubViewer(json.FileScrollOffsets, json.Epub.Files, json.CurrentPage, div);
            e.render();
        });
    }

    onMount(() => {
        let div = document.getElementById("book-view")!;
        utils.callApi("http://localhost:8080/cookie", "GET", {}).then((() => {
            getBook("Dune", div);
        }));
    });
</script>

{#if errorOut}
<div class="error">
    <p> Oops, something went wrong </p>
</div>
{:else}
<div class="container">
    <div class="left-sidepanel">
        <h1> {$book.Epub.Info.Title} </h1>
        <img alt="Ebook cover image" src={$book.Epub.CoverImagePath}/>
        <h3> {$book.Epub.Info.Author} </h3>
        <hr>
        <h5> {$book.Epub.Info.Description} </h5>
        <p> Date: {$book.Epub.Info.Date} </p>
        <p> Contributor: {$book.Epub.Info.Contributor} </p>
        <p> Coverage: {$book.Epub.Info.Coverage} </p>
        <p> Source: {$book.Epub.Info.Source} </p>
        <p> Rights: {$book.Epub.Info.Rights} </p>
        <p> Relation: {$book.Epub.Info.Relation} </p>
        <p> Publisher: {$book.Epub.Info.Publisher} </p>
        <p> Language: {$book.Epub.Info.Language} </p>
        <p> Identifier: {$book.Epub.Info.Identifier} </p>
        <p> Subjects: {#each $book.Epub.Info.Subjects as subject} {subject}  {/each} </p>
        <hr>
        <h3> Table of contents </h3>
        <ol>
            {#each $book.Epub.TableOfContents as section}
                <li><a href={section.Path}>{section.Name}</a></li>
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
