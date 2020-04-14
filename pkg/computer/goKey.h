
#include <windows.h>
#include <stddef.h>
#include <stdint.h>
#include <stdio.h>
#include <stdbool.h>

typedef int MMKeyCode;
#define DEADBEEF_MAX UINT32_MAX
static uint32_t deadbeef_seed;
static uint32_t deadbeef_beef = 0xdeadbeef;

uint32_t deadbeef_rand(void)
{
	deadbeef_seed = (deadbeef_seed << 7) ^ ((deadbeef_seed >> 25) + deadbeef_beef);
	deadbeef_beef = (deadbeef_beef << 7) ^ ((deadbeef_beef >> 25) + 0xdeadbeef);
	return deadbeef_seed;
}

#define DEADBEEF_RANDRANGE(a, b) \
	(uint32_t)DEADBEEF_UNIFORM(a, b)
#define WIN32_KEY_EVENT_WAIT(key, flags) \
    (win32KeyEvent(key, flags), Sleep(DEADBEEF_RANDRANGE(0, 1)))

#define H_INLINE static inline

enum _MMKeyFlags
{
    MOD_NONE = 0,
    MOD_META = MOD_WIN
};
typedef unsigned int MMKeyFlags;

/* Returns random double in the range [a, b).
 * Taken directly from the rand() man page. */
#define DEADBEEF_UNIFORM(a, b) \
	((a) + (deadbeef_rand() / (((double)DEADBEEF_MAX / (b - a) + 1))))
#define DEADBEEF_RANDRANGE(a, b) \
	(uint32_t)DEADBEEF_UNIFORM(a, b)
enum _MMKeyCode
{
    K_NOT_A_KEY = 9999,
    K_BACKSPACE = VK_BACK,
    K_DELETE = VK_DELETE,
    K_RETURN = VK_RETURN,
    K_TAB = VK_TAB,
    K_ESCAPE = VK_ESCAPE,
    K_UP = VK_UP,
    K_DOWN = VK_DOWN,
    K_RIGHT = VK_RIGHT,
    K_LEFT = VK_LEFT,
    K_HOME = VK_HOME,
    K_END = VK_END,
    K_PAGEUP = VK_PRIOR,
    K_PAGEDOWN = VK_NEXT,

    K_F1 = VK_F1,
    K_F2 = VK_F2,
    K_F3 = VK_F3,
    K_F4 = VK_F4,
    K_F5 = VK_F5,
    K_F6 = VK_F6,
    K_F7 = VK_F7,
    K_F8 = VK_F8,
    K_F9 = VK_F9,
    K_F10 = VK_F10,
    K_F11 = VK_F11,
    K_F12 = VK_F12,
    K_F13 = VK_F13,
    K_F14 = VK_F14,
    K_F15 = VK_F15,
    K_F16 = VK_F16,
    K_F17 = VK_F17,
    K_F18 = VK_F18,
    K_F19 = VK_F19,
    K_F20 = VK_F20,
    K_F21 = VK_F21,
    K_F22 = VK_F22,
    K_F23 = VK_F23,
    K_F24 = VK_F24,

    K_META = VK_LWIN,
    K_LMETA = VK_LWIN,
    K_RMETA = VK_RWIN,
    K_ALT = VK_MENU,
    K_LALT = VK_LMENU,
    K_RALT = VK_RMENU,
    K_CONTROL = VK_CONTROL,
    K_LCONTROL = VK_LCONTROL,
    K_RCONTROL = VK_RCONTROL,
    K_SHIFT = VK_SHIFT,
    K_LSHIFT = VK_LSHIFT,
    K_RSHIFT = VK_RSHIFT,
    K_CAPSLOCK = VK_CAPITAL,
    K_SPACE = VK_SPACE,
    K_PRINTSCREEN = VK_SNAPSHOT,
    K_INSERT = VK_INSERT,
    K_MENU = VK_APPS,

    K_NUMPAD_0 = VK_NUMPAD0,
    K_NUMPAD_1 = VK_NUMPAD1,
    K_NUMPAD_2 = VK_NUMPAD2,
    K_NUMPAD_3 = VK_NUMPAD3,
    K_NUMPAD_4 = VK_NUMPAD4,
    K_NUMPAD_5 = VK_NUMPAD5,
    K_NUMPAD_6 = VK_NUMPAD6,
    K_NUMPAD_7 = VK_NUMPAD7,
    K_NUMPAD_8 = VK_NUMPAD8,
    K_NUMPAD_9 = VK_NUMPAD9,
    K_NUMPAD_LOCK = VK_NUMLOCK,
    // VK_NUMPAD_
    K_NUMPAD_DECIMAL = VK_DECIMAL,
    K_NUMPAD_PLUS = VK_ADD,
    K_NUMPAD_MINUS = VK_SUBTRACT,
    K_NUMPAD_MUL = VK_MULTIPLY,
    K_NUMPAD_DIV = VK_DIVIDE,
    K_NUMPAD_CLEAR = K_NOT_A_KEY,
    K_NUMPAD_ENTER = VK_RETURN,
    K_NUMPAD_EQUAL = VK_OEM_PLUS,

    K_AUDIO_VOLUME_MUTE = VK_VOLUME_MUTE,
    K_AUDIO_VOLUME_DOWN = VK_VOLUME_DOWN,
    K_AUDIO_VOLUME_UP = VK_VOLUME_UP,
    K_AUDIO_PLAY = VK_MEDIA_PLAY_PAUSE,
    K_AUDIO_STOP = VK_MEDIA_STOP,
    K_AUDIO_PAUSE = VK_MEDIA_PLAY_PAUSE,
    K_AUDIO_PREV = VK_MEDIA_PREV_TRACK,
    K_AUDIO_NEXT = VK_MEDIA_NEXT_TRACK,
    K_AUDIO_REWIND = K_NOT_A_KEY,
    K_AUDIO_FORWARD = K_NOT_A_KEY,
    K_AUDIO_REPEAT = K_NOT_A_KEY,
    K_AUDIO_RANDOM = K_NOT_A_KEY,

    K_LIGHTS_MON_UP = K_NOT_A_KEY,
    K_LIGHTS_MON_DOWN = K_NOT_A_KEY,
    K_LIGHTS_KBD_TOGGLE = K_NOT_A_KEY,
    K_LIGHTS_KBD_UP = K_NOT_A_KEY,
    K_LIGHTS_KBD_DOWN = K_NOT_A_KEY
};

void win32KeyEvent(int key, MMKeyFlags flags){
	int scan = MapVirtualKey(key & 0xff, MAPVK_VK_TO_VSC);

	/* Set the scan code for extended keys */
	switch (key){
		case VK_RCONTROL:
		case VK_SNAPSHOT: /* Print Screen */
		case VK_RMENU: /* Right Alt / Alt Gr */
		case VK_PAUSE: /* Pause / Break */
		case VK_HOME:
		case VK_UP:
		case VK_PRIOR: /* Page up */
		case VK_LEFT:
		case VK_RIGHT:
		case VK_END:
		case VK_DOWN:
		case VK_NEXT: /* 'Page Down' */
		case VK_INSERT:
		case VK_DELETE:
		case VK_LWIN:
		case VK_RWIN:
		case VK_APPS: /* Application */
		case VK_VOLUME_MUTE:
		case VK_VOLUME_DOWN:
		case VK_VOLUME_UP:
		case VK_MEDIA_NEXT_TRACK:
		case VK_MEDIA_PREV_TRACK:
		case VK_MEDIA_STOP:
		case VK_MEDIA_PLAY_PAUSE:
		case VK_BROWSER_BACK:
		case VK_BROWSER_FORWARD:
		case VK_BROWSER_REFRESH:
		case VK_BROWSER_STOP:
		case VK_BROWSER_SEARCH:
		case VK_BROWSER_FAVORITES:
		case VK_BROWSER_HOME:
		case VK_LAUNCH_MAIL:
		{
			flags |= KEYEVENTF_EXTENDEDKEY;
			break;
		}
	}

	/* Set the scan code for keyup */
	// if ( flags & KEYEVENTF_KEYUP ) {
	// 	scan |= 0x80;
	// }

	// keybd_event(key, scan, flags, 0);
	
	INPUT keyInput;

	keyInput.type = INPUT_KEYBOARD;
	keyInput.ki.wVk = key;
	keyInput.ki.wScan = scan;
	keyInput.ki.dwFlags = flags;
	keyInput.ki.time = 0;
	keyInput.ki.dwExtraInfo = 0;
	SendInput(1, &keyInput, sizeof(keyInput));
}

struct KeyNames
{
    const char *name;
    MMKeyCode key;
} key_names[] = {
    {"backspace", K_BACKSPACE},
    {"delete", K_DELETE},
    {"enter", K_RETURN},
    {"tab", K_TAB},
    {"esc", K_ESCAPE},
    {"escape", K_ESCAPE},
    {"up", K_UP},
    {"down", K_DOWN},
    {"right", K_RIGHT},
    {"left", K_LEFT},
    {"home", K_HOME},
    {"end", K_END},
    {"pageup", K_PAGEUP},
    {"pagedown", K_PAGEDOWN},
    //
    {"f1", K_F1},
    {"f2", K_F2},
    {"f3", K_F3},
    {"f4", K_F4},
    {"f5", K_F5},
    {"f6", K_F6},
    {"f7", K_F7},
    {"f8", K_F8},
    {"f9", K_F9},
    {"f10", K_F10},
    {"f11", K_F11},
    {"f12", K_F12},
    {"f13", K_F13},
    {"f14", K_F14},
    {"f15", K_F15},
    {"f16", K_F16},
    {"f17", K_F17},
    {"f18", K_F18},
    {"f19", K_F19},
    {"f20", K_F20},
    {"f21", K_F21},
    {"f22", K_F22},
    {"f23", K_F23},
    {"f24", K_F24},
    //
    {"cmd", K_META},
    {"lcmd", K_LMETA},
    {"rcmd", K_RMETA},
    {"command", K_META},
    {"alt", K_ALT},
    {"lalt", K_LALT},
    {"ralt", K_RALT},
    {"ctrl", K_CONTROL},
    {"lctrl", K_LCONTROL},
    {"rctrl", K_RCONTROL},
    {"control", K_CONTROL},
    {"shift", K_SHIFT},
    {"lshift", K_LSHIFT},
    {"rshift", K_RSHIFT},
    {"right_shift", K_RSHIFT},
    {"capslock", K_CAPSLOCK},
    {"space", K_SPACE},
    {"print", K_PRINTSCREEN},
    {"printscreen", K_PRINTSCREEN},
    {"insert", K_INSERT},
    {"menu", K_MENU},

    {"audio_mute", K_AUDIO_VOLUME_MUTE},
    {"audio_vol_down", K_AUDIO_VOLUME_DOWN},
    {"audio_vol_up", K_AUDIO_VOLUME_UP},
    {"audio_play", K_AUDIO_PLAY},
    {"audio_stop", K_AUDIO_STOP},
    {"audio_pause", K_AUDIO_PAUSE},
    {"audio_prev", K_AUDIO_PREV},
    {"audio_next", K_AUDIO_NEXT},
    {"audio_rewind", K_AUDIO_REWIND},
    {"audio_forward", K_AUDIO_FORWARD},
    {"audio_repeat", K_AUDIO_REPEAT},
    {"audio_random", K_AUDIO_RANDOM},

    {"num0", K_NUMPAD_0},
    {"num1", K_NUMPAD_1},
    {"num2", K_NUMPAD_2},
    {"num3", K_NUMPAD_3},
    {"num4", K_NUMPAD_4},
    {"num5", K_NUMPAD_5},
    {"num6", K_NUMPAD_6},
    {"num7", K_NUMPAD_7},
    {"num8", K_NUMPAD_8},
    {"num9", K_NUMPAD_9},
    {"num_lock", K_NUMPAD_LOCK},

    {"num.", K_NUMPAD_DECIMAL},
    {"num+", K_NUMPAD_PLUS},
    {"num-", K_NUMPAD_MINUS},
    {"num*", K_NUMPAD_MUL},
    {"num/", K_NUMPAD_DIV},
    {"num_clear", K_NUMPAD_CLEAR},
    {"num_enter", K_NUMPAD_ENTER},
    {"num_equal", K_NUMPAD_EQUAL},

    {"numpad_0", K_NUMPAD_0},
    {"numpad_1", K_NUMPAD_1},
    {"numpad_2", K_NUMPAD_2},
    {"numpad_3", K_NUMPAD_3},
    {"numpad_4", K_NUMPAD_4},
    {"numpad_5", K_NUMPAD_5},
    {"numpad_6", K_NUMPAD_6},
    {"numpad_7", K_NUMPAD_7},
    {"numpad_8", K_NUMPAD_8},
    {"numpad_9", K_NUMPAD_9},
    {"numpad_lock", K_NUMPAD_LOCK},

    {"lights_mon_up", K_LIGHTS_MON_UP},
    {"lights_mon_down", K_LIGHTS_MON_DOWN},
    {"lights_kbd_toggle", K_LIGHTS_KBD_TOGGLE},
    {"lights_kbd_up", K_LIGHTS_KBD_UP},
    {"lights_kbd_down", K_LIGHTS_KBD_DOWN},

    {NULL, K_NOT_A_KEY} /* end marker */
};
MMKeyCode keyCodeForChar(const char c)
{
    return VkKeyScan(c);
}
H_INLINE void microsleep(double milliseconds)
{
    Sleep((DWORD)milliseconds);
}
int CheckKeyCodes(char *k, MMKeyCode *key)
{
    if (!key)
    {
        return -1;
    }
    if (strlen(k) == 1)
    {
        *key = keyCodeForChar(*k);
        return 0;
    }
    *key = K_NOT_A_KEY;
    struct KeyNames *kn = key_names;
    while (kn->name)
    {
    printf("\nk1=%s,k2=%s\n",k,kn->name);
        if (strcmp(k, kn->name) == 0)
        {
            *key = kn->key;
            break;
        }
        kn++;
    }
    if (*key == K_NOT_A_KEY)
    {
        return -2;
    }
    return 0;
}

int CheckKeyFlags(char *f, MMKeyFlags *flags)
{
    if (!flags)
    {
        return -1;
    }
    if (strcmp(f, "alt") == 0 || strcmp(f, "ralt") == 0 ||
        strcmp(f, "lalt") == 0)
    {
        *flags = MOD_ALT;
    }
    else if (strcmp(f, "command") == 0 || strcmp(f, "cmd") == 0 || strcmp(f, "rcmd") == 0 || strcmp(f, "lcmd") == 0)
    {
        *flags = MOD_META;
    }
    else if (strcmp(f, "control") == 0 || strcmp(f, "ctrl") == 0 ||
             strcmp(f, "rctrl") == 0 || strcmp(f, "lctrl") == 0)
    {
        *flags = MOD_CONTROL;
    }
    else if (strcmp(f, "shoft") == 0 || strcmp(f, "right_shift") == 0 ||
             strcmp(f, "rshift") == 0 || strcmp(f, "lashift") == 0)
    {
        *flags = MOD_SHIFT;
    }
    else if (strcmp(f, "none") == 0)
    {
        *flags = (MMKeyFlags)MOD_NONE;
    }
    else
    {
        return -2;
    }
    return 0;
}
int GetFlagsFromValue(char *value[], MMKeyFlags *flags, int num)
{
    if (!flags)
    {
        return -1;
    }
    int i;
    for (i = 0; i < num; i++)
    {
        MMKeyFlags f = MOD_NONE;
        const int rv = CheckKeyFlags(value[i], &f);
        if (rv)
        {
            return rv;
        }
        *flags = (MMKeyFlags)(*flags | f);
    }
    return 0;
}
void toggleKeyCode(MMKeyCode code, const bool down, MMKeyFlags flags)
{
    printf("\ncode:%d,down:%d,flags:%d\n",code,down,flags);
	const DWORD dwFlags = down ? 0 : KEYEVENTF_KEYUP;

	/* Parse modifier keys. */
	if (flags & MOD_META) WIN32_KEY_EVENT_WAIT(K_META, dwFlags);
	if (flags & MOD_ALT) WIN32_KEY_EVENT_WAIT(K_ALT, dwFlags);
	if (flags & MOD_CONTROL) WIN32_KEY_EVENT_WAIT(K_CONTROL, dwFlags);
	if (flags & MOD_SHIFT) WIN32_KEY_EVENT_WAIT(K_SHIFT, dwFlags);

	win32KeyEvent(code, dwFlags);
}
void tapKeyCode(MMKeyCode code, MMKeyFlags flags)
{
    toggleKeyCode(code, true, flags);
    toggleKeyCode(code, false, flags);
}
char *key_tap(char *k, char *akey, char *keyT, int keyDelay)
{
    printf("k=%s,akey=%s,keyT=%s,delay=%d",k,akey,keyT,keyDelay);
    MMKeyFlags flags = (MMKeyFlags)MOD_NONE;
    MMKeyCode key;
    if (strcmp(akey, "null") != 0)
    {
        if (strcmp(keyT, "null") == 0)
        {
            switch (CheckKeyFlags(akey, &flags))
            {
            case -1:
                return "Null pointer in key flag.";
                break;
            case -2:
                return "Invalid key flag specified.";
                break;
            }
        }
        else
        {
            char *akeyArr[2] = {akey, keyT};
            switch (GetFlagsFromValue(akeyArr, &flags, 2))
            {
            case -1:
                return "Null pointer in key flag.";
                break;
            case -2:
                return "Invalid key flag specified.";
                break;
            }
        }
    }
    switch (CheckKeyCodes(k, &key))
    {
    case -1:
        return "Null pointer in key code.";
        break;
    case -2:
        return "Invalid key code specified.";
        break;
    default:
        tapKeyCode(key, flags);
        microsleep(keyDelay);
    }

    return "";
}
char *key_Taps(char *k, char *keyArr[], int num, int keyDelay)
{
	MMKeyFlags flags = MOD_NONE;
	// MMKeyFlags flags = 0;
	MMKeyCode key;

	switch(GetFlagsFromValue(keyArr, &flags, num)) {
	// switch (CheckKeyFlags(akey, &flags)){
		case -1:
			return "Null pointer in key flag.";
			break;
		case -2:
			return "Invalid key flag specified.";
			break;
	}

	switch(CheckKeyCodes(k, &key)) {
		case -1:
			return "Null pointer in key code.";
			break;
		case -2:
			return "Invalid key code specified.";
			break;
		default:
			tapKeyCode(key, flags);
			microsleep(keyDelay);
	}

	// return "0";
	return "";
}