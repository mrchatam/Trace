import {
  createContext,
  useCallback,
  useContext,
  useEffect,
  useMemo,
  useState,
  type ReactNode,
} from 'react'
import { ApiError } from '../api/client'
import {
  getHealth,
  getProject,
  getVersion,
  type HealthResponse,
  type ProjectResponse,
  type VersionResponse,
} from '../api/ops'

const TOKEN_KEY = 'trace.gui.token'
const THEME_KEY = 'trace.gui.theme'

export type Theme = 'light' | 'dark'

type ChromeState = {
  theme: Theme
  token: string
  health: HealthResponse | null
  version: VersionResponse | null
  project: ProjectResponse | null
  chromeError: unknown
  loading: boolean
  storeReady: boolean
  setTheme: (t: Theme) => void
  setToken: (t: string) => void
  refresh: () => Promise<void>
}

const AppChromeContext = createContext<ChromeState | null>(null)

function readTheme(): Theme {
  const raw = localStorage.getItem(THEME_KEY)
  return raw === 'dark' ? 'dark' : 'light'
}

function applyTheme(theme: Theme) {
  document.documentElement.setAttribute('data-theme', theme)
}

export function AppChromeProvider({ children }: { children: ReactNode }) {
  const [theme, setThemeState] = useState<Theme>(() => readTheme())
  const [token, setTokenState] = useState(() => localStorage.getItem(TOKEN_KEY) ?? '')
  const [health, setHealth] = useState<HealthResponse | null>(null)
  const [version, setVersion] = useState<VersionResponse | null>(null)
  const [project, setProject] = useState<ProjectResponse | null>(null)
  const [chromeError, setChromeError] = useState<unknown>(null)
  const [loading, setLoading] = useState(true)

  useEffect(() => {
    applyTheme(theme)
  }, [theme])

  const setTheme = useCallback((t: Theme) => {
    localStorage.setItem(THEME_KEY, t)
    setThemeState(t)
  }, [])

  const setToken = useCallback((t: string) => {
    if (t) localStorage.setItem(TOKEN_KEY, t)
    else localStorage.removeItem(TOKEN_KEY)
    setTokenState(t)
  }, [])

  const refresh = useCallback(async () => {
    setLoading(true)
    setChromeError(null)
    const opt = { token: token || null }
    try {
      const [h, v] = await Promise.all([getHealth(opt), getVersion(opt)])
      setHealth(h)
      setVersion(v)
      try {
        const p = await getProject(opt)
        setProject(p)
      } catch (err) {
        setProject(null)
        setChromeError(err)
      }
    } catch (err) {
      setHealth(null)
      setVersion(null)
      setProject(null)
      setChromeError(err)
    } finally {
      setLoading(false)
    }
  }, [token])

  useEffect(() => {
    void refresh()
  }, [refresh])

  const storeReady = Boolean(project?.store_ready)

  const value = useMemo<ChromeState>(
    () => ({
      theme,
      token,
      health,
      version,
      project,
      chromeError,
      loading,
      storeReady,
      setTheme,
      setToken,
      refresh,
    }),
    [
      theme,
      token,
      health,
      version,
      project,
      chromeError,
      loading,
      storeReady,
      setTheme,
      setToken,
      refresh,
    ],
  )

  return <AppChromeContext.Provider value={value}>{children}</AppChromeContext.Provider>
}

export function useAppChrome(): ChromeState {
  const ctx = useContext(AppChromeContext)
  if (!ctx) throw new Error('useAppChrome requires AppChromeProvider')
  return ctx
}

export function useApiToken() {
  const { token } = useAppChrome()
  return token || null
}

export function unauthorizedHint(err: unknown): boolean {
  return err instanceof ApiError && err.code === 'UNAUTHORIZED'
}
