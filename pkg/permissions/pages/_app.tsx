import { Outlet } from "react-router-dom";
import RootLayout from "~/RootLayout";

export default function PermissionsLayout() {
  return (
    <RootLayout>
      <div class="box">
        <Outlet />
      </div>
    </RootLayout>
  )
}

export { action } from '.'
