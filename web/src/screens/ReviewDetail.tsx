import { useEffect, useState } from 'react'
import { Link, useParams } from 'react-router-dom'
import { ErrorBanner } from '../components/ErrorBanner'
import { useApiToken } from '../context/AppChrome'
import { getReview, type ReviewSummary } from '../api/ops'

export function ReviewDetail() {
  const { reviewId = '' } = useParams()
  const token = useApiToken()
  const [review, setReview] = useState<ReviewSummary | null>(null)
  const [error, setError] = useState<unknown>(null)

  async function load() {
    setError(null)
    try {
      setReview(await getReview(reviewId, { token }))
    } catch (err) {
      setError(err)
      setReview(null)
    }
  }

  useEffect(() => {
    if (reviewId) void load()
  }, [reviewId, token])

  return (
    <div>
      <p>
        <Link to="/reviews">← Reviews</Link>
      </p>
      <h1 className="page-title">{review?.title ?? 'Review'}</h1>

      {error ? <ErrorBanner error={error} onRetry={() => void load()} /> : null}

      {review ? (
        <div className="panel">
          <div className="stack">
            <div className="mono">{review.id}</div>
            {review.result ? <span className="pill">{review.result}</span> : null}
            {review.task_id ? (
              <div>
                Task:{' '}
                <Link to={`/tasks/${encodeURIComponent(review.task_id)}`}>{review.task_id}</Link>
              </div>
            ) : null}
            {review.body ? <p>{review.body}</p> : <p className="page-lead">No body.</p>}
          </div>
        </div>
      ) : null}
    </div>
  )
}
