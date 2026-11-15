import isRecord from "../core/IsRecord"
import { isString, isStringArray } from "../core/IsStringORStringArray"

export interface NetdiskSupportDetail {
    platform: string;
    features: Array<string>;
    main_file_path: string;
    type_file_path: string;
}

export interface ProcessedApiInfo {
    root_path: string;
    support_list: string[];
    support_details: Record<string, NetdiskSupportDetail>;
    version: string;
    license: string;
}

export default async function process_api_info(apiInfo: Record<string, unknown>): Promise<ProcessedApiInfo> {
    if (!apiInfo || !Object.hasOwn(apiInfo, "app")) {
        throw new Error("Api Info is empty or invalid");
    };
    if (apiInfo["app"] !== "MTD") {
        throw new Error("The key 'app' is not 'MTD'");
    };

    const root_path = Object.hasOwn(apiInfo, "root_path") ? apiInfo["root_path"] : null;
    if (root_path === null || root_path === undefined) {
        throw new Error("The key 'api_path' is not defined");
    } else if (!isString(root_path)) {
        throw new Error("The key 'root_path' is invalid")
    };

    const api_version = Object.hasOwn(apiInfo, "version") ? apiInfo["version"] : "null";
    if (!isString(api_version)) {
        throw new Error("The key 'version' is invalid");
    };
    const api_license = Object.hasOwn(apiInfo, "license") ? apiInfo["license"] : "null";
    if (!isString(api_license)) {
        throw new Error("The key 'license' is invalid");
    };

    if (!Object.hasOwn(apiInfo, "netdisk_support")) {
        throw new Error("The key 'netdisk_support' is not defined");
    } else if (!isStringArray(apiInfo["netdisk_support"]) || apiInfo["netdisk_support"].length === 0) {
        throw new Error("The key 'netdisk_support' is not a string array or empty");
    } else {
        const netdisk_support_list: string[] = apiInfo["netdisk_support"];
        const support_details: Record<string, NetdiskSupportDetail> = Object.create(null);
        for (const key of netdisk_support_list) {
            const each_platform_details = apiInfo[`${key}_support`];

            if (!each_platform_details || !Object.hasOwn(each_platform_details, "platform")) {
                throw new Error(`The key '${key}_support' is empty or invalid`);
            } else if (!isRecord(each_platform_details)) {
                throw new Error(`The key '${key}_support' is invalid`);
            };

            const platform = Object.hasOwn(each_platform_details, "platform") ? each_platform_details["platform"] : null;
            if (!isString(platform)) {
                throw new Error(`The key 'platform' in the key '${key}_support' is invalid`);
            };

            

            const main_file_path = Object.hasOwn(each_platform_details, "main_file_path") ? each_platform_details["main_file_path"] : null;
            if (!isString(main_file_path)) {
                throw new Error(`The key 'main_file_path' in the key '${key}_support' is invalid`);
            };

            const type_file_path = Object.hasOwn(each_platform_details, "type_file_path") ? each_platform_details["type_file_path"] : null;
            if (!isString(type_file_path)) {
                throw new Error(`The key 'type_file_path' in the key '${key}_support' is invalid`);
            };

            const features = Object.hasOwn(each_platform_details, "features") ? each_platform_details["features"] : null;
            if (!isStringArray(features)) {
                throw new Error(`The key 'features' in the key '${key}_support' is invalid`);
            };
            if (features.length === 0) {
                throw new Error(`The key 'features' in the key '${key}_support' is empty`);
            };

            if (platform in support_details) {
                throw new Error(`The key '${key}' in the netdisk_support_list is duplicated`);
                // console.log(`The key '${key}' in the netdisk_support_list is duplicated`);
                // continue;
            } else {
                const details: NetdiskSupportDetail = {"platform": platform, "features": features, "main_file_path": main_file_path, "type_file_path": type_file_path}; // return {"platform1": {details1...}, "platform2": {details2...}}
                support_details[platform] = details;
            };
        };
        return {"root_path": root_path, "support_list": netdisk_support_list, "support_details": support_details, "version": api_version, "license": api_license};
    };
}