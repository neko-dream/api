import type { BaseRecord, DataProvider } from "@refinedev/core";
import type { AxiosInstance } from "axios";
import type { ManagedUser } from "../pages/users/types.ts";

// Only adapt the operator user API. Existing resources retain their provider.
export function withUserManagement(
  client: AxiosInstance,
  fallback: DataProvider,
): DataProvider {
  return {
    ...fallback,
    create: async <
      TData extends BaseRecord = BaseRecord,
      TVariables = Record<string, unknown>,
    >(
      params: Parameters<DataProvider["create"]>[0] & { variables: TVariables },
    ) => {
      if (params.resource !== "users")
        return fallback.create<TData, TVariables>(params);
      const { data } = await client.post<{
        userID: string;
        invitationSent: boolean;
      }>("/v1/manage/users", params.variables);
      return { data: { ...data, id: data.userID } as unknown as TData };
    },
    getList: async <TData extends BaseRecord = BaseRecord>(
      params: Parameters<DataProvider["getList"]>[0],
    ) => {
      if (params.resource !== "users") return fallback.getList<TData>(params);
      const page = params.pagination?.currentPage ?? 1;
      const limit = params.pagination?.pageSize ?? 20;
      const filter = params.filters?.find(
        (item) => "field" in item && item.field === "search",
      );
      const search =
        filter && "value" in filter ? String(filter.value ?? "") : "";
      const { data } = await client.get<{
        users: Omit<ManagedUser, "id">[];
        totalCount: number;
      }>("/v1/manage/users/list", {
        params: { offset: (page - 1) * limit, limit, search },
      });
      return {
        data: data.users.map((user) => ({
          ...user,
          id: user.userID,
        })) as unknown as TData[],
        total: data.totalCount,
      };
    },
    getOne: async <TData extends BaseRecord = BaseRecord>(
      params: Parameters<DataProvider["getOne"]>[0],
    ) => {
      if (params.resource !== "users") return fallback.getOne<TData>(params);
      const { data } = await client.get<Omit<ManagedUser, "id">>(
        `/v1/manage/users/${encodeURIComponent(String(params.id))}`,
      );
      return { data: { ...data, id: data.userID } as unknown as TData };
    },
  };
}
