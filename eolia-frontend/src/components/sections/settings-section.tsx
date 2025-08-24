import { useState } from "react";
import Toggle from "../ui/toggle";

const SettingsSection = () => {
  const [dailyLimitEnabled, setDailyLimitEnabled] = useState(true);
  const [dailyLimitValue, setDailyLimitValue] = useState("250");

  const [txApprovalRequired, setTxApprovalRequired] = useState(false);
  const [notificationsEnabled, setNotificationsEnabled] = useState(true);
  const [autoSignSmallTx, setAutoSignSmallTx] = useState(false);
  const [autoSignLimit, setAutoSignLimit] = useState("10");

  const [sessionTimeout, setSessionTimeout] = useState("30");

  return (
    <div className="w-[400px] h-[640px] relative flex flex-col items-center rounded-t-xl bg-[#f7f7f7] px-[40px] overflow-y-scroll scroll-smooth py-8">
      <h2 className="text-xl font-bold text-gray-800 mb-6">Settings</h2>

      {/* Daily Spending Limit */}
      <div className="w-full mb-6">
        <div className="flex justify-between items-center">
          <p className="text-[14px] font-medium">Daily Spending Limit</p>
          <Toggle isOn={dailyLimitEnabled} onToggle={() => setDailyLimitEnabled(!dailyLimitEnabled)} />
        </div>
        {dailyLimitEnabled && (
          <div className="relative">
            <p className="absolute text-sm text-gray-500 top-1/2 -translate-y-1/3 right-4">$</p>
            <input
              type="price"
              value={dailyLimitValue}
              onChange={(e) => setDailyLimitValue(e.target.value)}
              className="mt-2 w-full px-3 py-2 rounded-md bg-white text-sm border border-gray-300 focus:outline-none"
              placeholder="Enter USD limit"
            />
          </div>
        )}
      </div>

      {/* Manual Transaction Approval */}
      <div className="w-full mb-6">
        <div className="flex justify-between items-center">
          <p className="text-[14px] font-medium">Manual Tx Approval</p>
          <Toggle isOn={txApprovalRequired} onToggle={() => setTxApprovalRequired(!txApprovalRequired)} />
        </div>
        <p className="text-xs text-gray-500 mt-1">
          Require manual approval before each transaction.
        </p>
      </div>

      {/* Push Notifications */}
      <div className="w-full mb-6">
        <div className="flex justify-between items-center">
          <p className="text-[14px] font-medium">Push Notifications</p>
          <Toggle isOn={notificationsEnabled} onToggle={() => setNotificationsEnabled(!notificationsEnabled)} />
        </div>
        <p className="text-xs text-gray-500 mt-1">
          Receive updates on status and account activity.
        </p>
      </div>

      {/* Auto Sign Small Tx */}
      <div className="w-full mb-6">
        <div className="flex justify-between items-center">
          <p className="text-[14px] font-medium">Auto Sign Small Tx</p>
          <Toggle isOn={autoSignSmallTx} onToggle={() => setAutoSignSmallTx(!autoSignSmallTx)} />
        </div>
        {autoSignSmallTx && (
          <input
            type="number"
            value={autoSignLimit}
            onChange={(e) => setAutoSignLimit(e.target.value)}
            className="mt-2 w-full px-3 py-2 rounded-md bg-white text-sm border border-gray-300 focus:outline-none"
            placeholder="Max $ amount"
          />
        )}
        <p className="text-xs text-gray-500 mt-1">
          Auto-sign transactions under this amount.
        </p>
      </div>

      {/* Session Key Timeout */}
      <div className="w-full mb-6">
        <p className="text-[14px] font-medium mb-2">Session Timeout (minutes)</p>
        <input
          type="number"
          value={sessionTimeout}
          onChange={(e) => setSessionTimeout(e.target.value)}
          className="w-full px-3 py-2 rounded-md bg-white text-sm border border-gray-300 focus:outline-none"
          placeholder="Timeout in minutes"
        />
        <p className="text-xs text-gray-500 mt-1">
          After this period, session key will expire automatically.
        </p>
      </div>
    </div>
  );
};

export default SettingsSection;