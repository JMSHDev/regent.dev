import { authRequest } from "./auth";

interface Device {
  name: string;
  customer: string;
  status: string;
  lastUpdated: string;
  activated: boolean;
}

const getDeviceList = async () => {
  const apiResp = await authRequest.get("/api/devices/?format=json");
  const deviceList: Device[] = [];

  for (const device of apiResp.data) {
    deviceList.push({
      name: device.name,
      customer: device.customer,
      status: device.status,
      lastUpdated: device.last_updated,
      activated: device.auth[0] === "activated"
    });
  }

  return deviceList;
};

export { getDeviceList, Device };
