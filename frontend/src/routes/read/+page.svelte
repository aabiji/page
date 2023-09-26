<script lang="ts">
    import BookViewer from "./bookview.svelte";
    import Navbar from "../navbar.svelte";
    import { onMount } from "svelte";
    import * as utils from "$lib/utils";

    onMount(() => {
        utils.redirectIfNotAuth();
    });
</script>

<Navbar />
<BookViewer bookName="Dune" />
