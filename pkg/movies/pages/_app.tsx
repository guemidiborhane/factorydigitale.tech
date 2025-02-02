import { Outlet } from "react-router-dom";
import RootLayout from "~/RootLayout";

export default function MoviesLayout() {
  return (
    <RootLayout>
      <Outlet />
    </RootLayout>
  )
}
