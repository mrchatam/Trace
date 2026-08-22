import { BrowserRouter, Navigate, Route, Routes } from 'react-router-dom'
import { AppChromeProvider } from './context/AppChrome'
import { Shell } from './layout/Shell'
import { Overview } from './screens/Overview'
import { Tasks } from './screens/Tasks'
import { TaskDetail } from './screens/TaskDetail'
import { Loop } from './screens/Loop'
import { Graph } from './screens/Graph'
import { Discoveries } from './screens/Discoveries'
import { Reviews } from './screens/Reviews'
import { ReviewDetail } from './screens/ReviewDetail'
import { Seed } from './screens/Seed'
import { Settings } from './screens/Settings'

export default function App() {
  return (
    <AppChromeProvider>
      <BrowserRouter>
        <Routes>
          <Route element={<Shell />}>
            <Route index element={<Graph />} />
            <Route path="overview" element={<Overview />} />
            <Route path="graph" element={<Navigate to="/" replace />} />
            <Route path="tasks" element={<Tasks />} />
            <Route path="tasks/:taskId" element={<TaskDetail />} />
            <Route path="loop" element={<Loop />} />
            <Route path="discoveries" element={<Discoveries />} />
            <Route path="reviews" element={<Reviews />} />
            <Route path="reviews/:reviewId" element={<ReviewDetail />} />
            <Route path="seed" element={<Seed />} />
            <Route path="settings" element={<Settings />} />
            <Route path="*" element={<Navigate to="/" replace />} />
          </Route>
        </Routes>
      </BrowserRouter>
    </AppChromeProvider>
  )
}
