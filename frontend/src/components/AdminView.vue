<template>
    <div class="container">
        <h1>UniLib - Панель администратора</h1>

        <nav class="admin-nav">
            <router-link to="/admin/manage-books" class="nav-link">Управление книгами</router-link>
            <router-link to="/admin/manage-readers" class="nav-link">Управление читателями</router-link>
            <router-link to="/admin/reports" class="nav-link">Отчеты</router-link>
            <router-link to="/admin/report-history" class="nav-link">История отчетов</router-link>
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
.admin-nav {
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