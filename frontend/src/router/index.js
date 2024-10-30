import { createRouter, createWebHistory } from 'vue-router';
import HomeView from '../views/HomeView.vue';
import RegisterView from '../views/RegisterView.vue';
import AdminView from '../components/AdminView.vue';
import LibrarianView from '../components/LibrarianView.vue';
import ReaderView from '../components/ReaderView.vue';
import ReaderHistoryView from '../components/ReaderHistoryView.vue';
import LibrarianRegisterReaderView from '../components/LibrarianRegisterReaderView.vue';
import LibrarianIssueBookView from '../components/LibrarianIssueBookView.vue';
import LibrarianReturnBookView from '../components/LibrarianReturnBookView.vue';
import LibrarianSearchBookView from '../components/LibrarianSearchBookView.vue';
import AdminManageBooksView from '../components/AdminManageBooksView.vue';
import AdminManageReadersView from '../components/AdminManageReadersView.vue';
import AdminReportsView from '../components/AdminReportsView.vue';
import AdminReportHistoryView from '../components/AdminReportHistoryView.vue';


const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes: [
    {
      path: '/',
      name: 'home',
      component: HomeView
    },
    {
      path: '/register',
      name: 'register',
      component: RegisterView
    },
    {
      path: '/admin',
      name: 'admin',
      component: AdminView,
      children: [
        {
          path: 'manage-books',
          name: 'adminManageBooks',
          component: AdminManageBooksView
        },
        {
          path: 'manage-readers',
          name: 'adminManageReaders',
          component: AdminManageReadersView
        },
        {
          path: 'reports',
          name: 'adminReports',
          component: AdminReportsView
        },
        {
          path: '/admin/report-history',
          name: 'adminReportHistory',
          component: AdminReportHistoryView
        }
      ]
    },
    {
      path: '/librarian',
      name: 'librarian',
      component: LibrarianView,
      children: [
        {
          path: 'register-reader',
          name: 'librarianRegisterReader',
          component: LibrarianRegisterReaderView
        },
        {
          path: 'issue-book',
          name: 'librarianIssueBook',
          component: LibrarianIssueBookView
        },
        {
          path: 'return-book',
          name: 'librarianReturnBook',
          component: LibrarianReturnBookView
        },
        {
          path: 'search-book',
          name: 'librarianSearchBook',
          component: LibrarianSearchBookView
        }
      ]
    },
    {
      path: '/reader',
      name: 'reader',
      component: ReaderView
    },
    {
      path: '/reader/history',
      name: 'readerHistory',
      component: ReaderHistoryView
    }
  ]
});

export default router;