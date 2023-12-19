/** This is an automatically generated class by FairyGUI. Please do not modify it. **/

import UI_TopUI from "./UI_TopUI";
import UI_BottomUI from "./UI_BottomUI";
import UI_MainUI from "./UI_MainUI";

export default class MainUIBinder {
	public static bindAll():void {
		fgui.UIObjectFactory.setExtension(UI_TopUI.URL, UI_TopUI);
		fgui.UIObjectFactory.setExtension(UI_BottomUI.URL, UI_BottomUI);
		fgui.UIObjectFactory.setExtension(UI_MainUI.URL, UI_MainUI);
	}
}