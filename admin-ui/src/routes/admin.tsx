import { createFileRoute, Navigate } from "@tanstack/react-router";

export const Route = createFileRoute('/admin')({
  component: AdminRedirect,
})


function AdminRedirect() {
  return <Navigate to="/" />
}
