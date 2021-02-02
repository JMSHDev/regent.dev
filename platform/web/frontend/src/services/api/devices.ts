import { authRequest } from "./auth";

interface Device {
  name: string;
  customer: string;
  status: string;
  lastUpdated: string;
  activated: boolean;
}

const getDeviceList = async () => {
  const apiResp = await authRequest.get("/api/devices/");
  const deviceList: Device[] = [];

  for (const device of apiResp.data.results) {
    deviceList.push({
      name: device.name,
      customer: device.customer,
      status: device.status,
      lastUpdated: new Date(device.last_updated).toLocaleString("en-GB", { timeZone: "UTC" }),
      activated: device.auth[0] === "activated",
    });
  }

  return deviceList;
};

export { getDeviceList, Device };
