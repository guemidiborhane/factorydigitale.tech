import { Form, type ActionFunctionArgs, type LoaderFunctionArgs, redirect } from "react-router-dom";
import { parseForm } from "~/helpers/form";
import { useLoaderData } from "~/hooks";
import { fetchApi } from "~/helpers/http";
import Card from "@/permissions/components/Card";
import { useCan } from "~/hooks/useCan";
import type { FC } from "preact/compat";
import Modal from "ui/Modal";
import Button from "ui/Button";
import { APISchemas } from "~/api";
import { useT } from "i18n/index";
import { batch, useSignal } from "@preact/signals";

const Inputs: FC<{ role: string }> = ({ role }) => {
  const response = useLoaderData<'/api/permissions'>()

  if (!response.ok) return null

  const { permissions, roles_permissions } = response.data

  return (
    <div class="grid grid-cols-1 gap-1 md:grid-cols-3">
      {permissions && Object.entries(permissions).map(([resource, perms]) => (
        <div class="box">
          <h6 class="text-center">{resource}</h6>
          {perms.map((action: string, index: number) => {
            const key = `permissions:${resource}:${index}`
            const checked = roles_permissions && roles_permissions[role][resource]?.indexOf(action) > -1
            return (
              <label for={key} class="flex justify-between">
                {action}
                <input type="checkbox" name={key} class="ms-1" checked={checked} value={action} />
              </label>
            )
          })}
        </div>
      ))}
    </div>
  )
}

export async function action({ request }: ActionFunctionArgs) {
  const { signal } = request
  const { permissions: body, role } = await parseForm(request) as { permissions: APISchemas['pkg_permissions.PermissionsParams'], role: string }

  return await fetchApi('/api/permissions/{role}', { method: 'put', body, signal, urlParams: { role } })
}

export async function loader({ request }: LoaderFunctionArgs) {
  const response = await fetchApi('/api/permissions', { signal: request.signal })
  if (!response.ok) return redirect('/')
  return response
}

export default function IndexPermissions() {
  const response = useLoaderData<'/api/permissions'>()

  if (!response.ok) return null

  const { roles_permissions } = response.data
  const show = useSignal<boolean>(false)
  const role = useSignal<string | undefined>(undefined)
  const canEdit = useCan("permissions:update")
  const { t } = useT()

  return (
    <>
      <div class="flex flex-col gap-3 justify-center items-start lg:flex-row">
        {Object.entries(roles_permissions).map(([r]) => {
          return r != 'root' && (
            <div class="py-2 px-4 cursor-pointer hover:bg-gray-100" onClick={() => {
              batch(() => {
                role.value = r
                show.value = true
              })
            }}>
              <Card role={r} permissions={roles_permissions[r]} />
            </div>
          )
        })}
      </div >
      {role && canEdit && (
        <Modal show={show}>
          <Form method="PUT" class="flex flex-col justify-between min-h-full" dir="ltr">
            <input type="hidden" name="role" value={role} />
            {role.value && <Inputs role={role.value} />}
            <Button type="Success" class="col-span-3">{t('misc.submit')}</Button>
          </Form >
        </Modal>
      )
      }
    </>
  )
}
