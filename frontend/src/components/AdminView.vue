<template>
    <div class="container">
        <h1>UniLib - Панель администратора</h1>

        <nav>
            <router-link to="/admin/manage-books">Управление книгами</router-link>
            <router-link to="/admin/manage-readers">Управление читателями</router-link>
            <router-link to="/admin/reports">Отчеты</router-link>
        </nav>

        <router-view />

        <button @click="logout">Выход</button>
    </div>
</template>

<script>
import { useRouter } from 'vue-router';

export default {
    setup() {
        const router = useRouter();

        const logout = () => {
            //  Логика  выхода 
            router.push('/');
        };

        return { logout };
    },
};
</script>

<style scoped>
</style>