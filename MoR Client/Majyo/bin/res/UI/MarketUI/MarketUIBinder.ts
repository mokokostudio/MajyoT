/** This is an automatically generated class by FairyGUI. Please do not modify it. **/

import UI_MarketUI from "./UI_MarketUI";
import UI_Market_Info from "./UI_Market_Info";

export default class MarketUIBinder {
	public static bindAll():void {
		fgui.UIObjectFactory.setExtension(UI_MarketUI.URL, UI_MarketUI);
		fgui.UIObjectFactory.setExtension(UI_Market_Info.URL, UI_Market_Info);
	}
}