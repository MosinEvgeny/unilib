<template>
  <div>
    <h1>UniLib</h1>
    <ul>
      <li v-for="book in books" :key="book.book_id">
        {{ book.title }} - {{ book.author }}
      </li>
    </ul>
  </div>
</template>

<script>
import { onMounted, ref } from 'vue';

export default {
  setup() {
    const books = ref([]);

    onMounted(async () => {
      try {
        const response = await fetch('http://localhost:8080/books');
        books.value = await response.json();
      } catch (error) {
        console.error('Ошибка при загрузке данных:', error);
      }
    });

    return { books };
  },
};
</script>