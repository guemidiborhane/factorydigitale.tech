import { type LoaderFunctionArgs, redirect } from "react-router-dom";
import { useLoaderData } from "~/hooks";
import { fetchApi } from "~/helpers/http";
import { ListMovies } from "@/movies/components";

export async function loader({ request }: LoaderFunctionArgs) {
  const response = await fetchApi('/api/movies/favourites', { signal: request.signal })
  if (!response.ok) return redirect('/')
  return response
}

export default function IndexMovies() {
  const response = useLoaderData<'/api/movies'>()

  if (!response.ok) return null

  const movies = response.data

  return (
    <>
      <ListMovies movies={movies} />
    </>
  )
}
