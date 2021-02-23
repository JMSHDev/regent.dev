import { authRequest } from "./auth";

interface Device {
  name: string;
  customer: string;
  status: string;
  lastUpdated: string;
  activated: boolean;
  pk: string;
}

const nullDevice: Device = {
  name: "",
  customer: "",
  status: "",
  lastUpdated: "",
  activated: false,
  pk: ""
}

const getDeviceList = async () => {
  const apiResp = await authRequest.get("/api/devices/");
  const deviceList: Device[] = [];

  for (const device of apiResp.data.results) {
    const deviceUrlComponents = device.url.split("/");

    deviceList.push({
      name: device.name,
      customer: device.customer,
      status: device.status,
      lastUpdated: new Date(device.last_updated).toLocaleString("en-GB", { timeZone: "UTC" }),
      activated: device.auth[0] === "activated",
      pk: deviceUrlComponents[deviceUrlComponents.length - 2]
    });
  }

  return deviceList;
};

const getDevice = async (pk: string) => {
  const apiResp = await authRequest.get(`/api/devices/${pk}/`);
  const apiDevice = apiResp.data;

  const device = {
      name: apiDevice.name,
      customer: apiDevice.customer,
      status: apiDevice.status,
      lastUpdated: new Date(apiDevice.last_updated).toLocaleString("en-GB", { timeZone: "UTC" }),
      activated: apiDevice.auth[0] === "activated",
      pk: pk
    }

  return device;
};

export { getDeviceList, getDevice, nullDevice, Device };
