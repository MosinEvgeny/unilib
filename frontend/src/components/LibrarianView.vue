<template>
    <div class="container">
        <h1>UniLib - Панель библиотекаря</h1>

        <nav class="librarian-nav">
            <router-link to="/librarian/search-book" class="nav-link">Поиск книг</router-link>
            <router-link to="/librarian/register-reader" class="nav-link">Регистрация читателя</router-link>
            <router-link to="/librarian/issue-book" class="nav-link">Выдача книги</router-link>
            <router-link to="/librarian/return-book" class="nav-link">Прием книги</router-link>
            <button @click="logout">Выход</button>
        </nav>

        <router-view />

    </div>
</template>

<script>
import { useRouter } from 'vue-router';

export default {
    setup() {
        const router = useRouter();

        const logout = () => {
            router.push('/');
        };

        return { logout };
    },
};
</script>

<style scoped>
.librarian-nav {
    display: flex;
    gap: 1rem;
    margin-bottom: 2rem;
}

.nav-link {
    padding: 0.8rem 1.6rem;
    border-radius: 4px;
    background-color: #fbd604;
    color: #fff;
    text-decoration: none;
}

.nav-link:hover {
    background-color: darken(#fbd604, 10%);
}

.router-link-exact-active {
    background-color: #e6330f;
}
</style>