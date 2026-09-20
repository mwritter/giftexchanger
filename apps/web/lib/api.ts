export type StatusResponse = {
  status: string;
};

export async function fetchApiStatus(
  path: "/api/health" | "/api/ready",
): Promise<StatusResponse> {
  const response = await fetch(path);
  const data = (await response.json()) as StatusResponse;

  if (!response.ok) {
    throw new Error(data.status ?? response.statusText);
  }

  return data;
}
