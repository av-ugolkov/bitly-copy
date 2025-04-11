import { createBrowserRouter } from 'react-router';

import BaseLayout from './layouts/Root';
import Main from './pages/Main';
import Statistics from './pages/Statistics';
import UserUrls from './pages/UserUrls';

export const router = createBrowserRouter([
  {
    path: '/',
    element: <BaseLayout />,
    children: [
      {
        path: '/',
        element: <Main />,
      },
      {
        path: '/get-urls',
        element: <UserUrls />,
      },
      {
        path: '/statistics',
        element: <Statistics />,
      },
    ],
  },
]);
