#include <stdio.h>
#include <windows.h>
typedef struct
{
    HWND hWnd;
    DWORD dwPid;
} WNDINFO;

BOOL CALLBACK EnumWindowsProc2(HWND hWnd, LPARAM lParam)
{
    WNDINFO *pInfo = (WNDINFO *)lParam;
    DWORD dwProcessId = 0;
    GetWindowThreadProcessId(hWnd, &dwProcessId);

    if (dwProcessId == pInfo->dwPid)
    {
        pInfo->hWnd = hWnd;
        return FALSE;
    }
    return TRUE;
}

HWND GetHwndByPid(DWORD dwProcessId)
{
    WNDINFO info = {0};
    info.hWnd = NULL;
    info.dwPid = dwProcessId;
    EnumWindows(EnumWindowsProc2, (LPARAM)&info);
    return info.hWnd;
}
void ActiveForcePid(int pid)
{
    DWORD dwTimeout = -1;
    SystemParametersInfo(SPI_GETFOREGROUNDLOCKTIMEOUT, 0, (LPVOID)&dwTimeout, 0);
    if (dwTimeout >= 100)
    {
        SystemParametersInfo(SPI_SETFOREGROUNDLOCKTIMEOUT, 0, (LPVOID)0, SPIF_SENDCHANGE | SPIF_UPDATEINIFILE);
    }
    HWND hForeWnd = NULL;
    HWND hWnd = GetHwndByPid(pid);
    DWORD dwForeID;
    DWORD dwCurID;
    hForeWnd = GetForegroundWindow();
    dwCurID = GetCurrentThreadId();
    dwForeID = GetWindowThreadProcessId(hForeWnd, NULL);
    AttachThreadInput(dwCurID, dwForeID, TRUE);

    //    ShowWindow(hWnd, SW_SHOWNORMAL);
    SetWindowPos(hWnd, HWND_TOPMOST, 0, 0, 0, 0, SWP_NOSIZE | SWP_NOMOVE);
    SetWindowPos(hWnd, HWND_NOTOPMOST, 0, 0, 0, 0, SWP_NOSIZE | SWP_NOMOVE);
    SetForegroundWindow(hWnd);
    //    SetFocus( hWnd );
    //    SendMessage( hWnd, WM_SETFOCUS, 0, 0 );
    //    PostMessage( hWnd, WM_SETFOCUS, 0, 0 );
    //    SendMessage( hWnd, WM_APP, 0, 0 );
    //    PostMessage( hWnd, WM_APP, 0, 0 );

    SetFocus(hWnd);
    SetWindowPos(hWnd,                     //编辑框  窗口句柄。
                 HWND_TOP,                 //将指定的内容置于Z顺序的顶部。
                 0,                        // x-没关系
                 0,                        // y-没关系
                 0,                        //宽度-没关系
                 0,                        //高度-没关系
                 SWP_NOSIZE | SWP_NOMOVE); //忽略x，y，width和height参数。
    AttachThreadInput(dwCurID, dwForeID, FALSE);
}

void ClickWindowByPid(int pid)
{
    HWND hwnd = NULL;
    hwnd = GetHwndByPid(pid);
    if (hwnd == NULL)
    {
        return;
    }
    RECT rect;
    GetWindowRect(hwnd, &rect);
    int err=GetLastError();
    printf("w:%d,h:%d,err:%d\n",rect.left,rect.right,err);
}