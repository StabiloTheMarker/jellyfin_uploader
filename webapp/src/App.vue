<script setup lang="ts">

import {Button} from "@/components/ui/button";
</script>

<template>
  <div id="container" class="h-screen w-screen flex justify-center items-center">
    <form action="/upload"  enctype="multipart/form-data" method="post">

      <input name="files" type="file">
      <Button type="submit">Submit</Button>
    </form>
  </div>
</template>