<script>
    import { onMount } from "svelte";
    import * as utils from "$lib/utils";
    import Navbar from "./navbar.svelte";
    import Book from "./book.svelte";

    onMount(() => {
        utils.redirectIfNotAuth();
    });
</script>

<Navbar />
<h1> Home page </h1>
<Book name="Book namejfoieajfeoiwj" userBookId="1iue092ija"/>
