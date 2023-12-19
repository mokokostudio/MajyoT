/** This is an automatically generated class by FairyGUI. Please do not modify it. **/

import UI_PVPUI from "./UI_PVPUI";
import UI_PVP_Info from "./UI_PVP_Info";

export default class PVPUIBinder {
	public static bindAll():void {
		fgui.UIObjectFactory.setExtension(UI_PVPUI.URL, UI_PVPUI);
		fgui.UIObjectFactory.setExtension(UI_PVP_Info.URL, UI_PVP_Info);
	}
}