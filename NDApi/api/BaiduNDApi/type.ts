export interface FileListItem {
    name: string;
    path: string;
    isdir: number; // 1 for folder, 0 for file
    fid: number;
    size: number;
}

export interface FileInfo {
  fs_id: number;
  filename: string;
  size: number;
  isdir: number; // 1 for folder, 0 for file
  dlink: string | null;
}
