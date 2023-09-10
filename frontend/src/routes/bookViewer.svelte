<script lang="ts">
    import { onMount } from "svelte";
    import { EpubViewer } from "./epub";
    import { callApi } from "./utils";

    function getBookInfo(name: string) {
        return callApi(`http://localhost:8080/book/get/${name}`, "GET", {});
    }

    let book = {Title: "AnimalFarm", Author: "-", Date: "-", Description: "-",
                Contributor: "-", Coverage: "-", Source: "-", Rights: "-",
                Relation: "-", Publisher: "-", Language: "-", Identifier: "-", Subjects: []};
    onMount(() => {
        let div = document.getElementById("book-view")!;

        getBookInfo(book.Title).then((json) => {
            if ("Server error" in json) return;

            book = json.Info;
            book.Subjects = book.Subjects == null ? [] : book.Subjects;

            let e = new EpubViewer(json.Files, div);
            e.render();
        });
    });
</script>

<div class="container">
    <div class="info-sidepanel">
        <h1> {book.Title} </h1>
        <!--- cover image goes here: --->
        <h3> {book.Author} </h3>
        <h5> {book.Description} </h5>
        <p> Date: {book.Date} </p>
        <p> Contributor: {book.Contributor} </p>
        <p> Coverage: {book.Coverage} </p>
        <p> Source: {book.Source} </p>
        <p> Rights: {book.Rights} </p>
        <p> Relation: {book.Relation} </p>
        <p> Publisher: {book.Publisher} </p>
        <p> Language: {book.Language} </p>
        <p> Identifier: {book.Identifier} </p>
        <p> Subjects: {#each book.Subjects as subject} {subject}, {/each} </p>
        <hr>
    </div>
    <div id="book-view"></div>
</div>

<style>
    .container {
        display: flex;
        padding: none;
    }

    #book-view {
        width: 600px;
        margin: 0 auto;
    }

    .info-sidepanel {
        width: 25%;
        height: 100vh;
        text-align: center;
        background-color: blue;
    }
</style>
