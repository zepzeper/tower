// dashboard/src/components/Header.jsx
import { useState } from 'react';
import { Bell, Search, Menu } from 'lucide-react';

const TopBar = () => {
  return (
    <header className="bg-white shadow h-16 flex items-center px-6">
      <button className="lg:hidden mr-4">
        <Menu className="h-6 w-6 text-gray-500" />
      </button>
      
      <div className="relative flex-1 max-w-xs">
        <div className="absolute inset-y-0 left-0 flex items-center pl-3">
          <Search className="h-5 w-5 text-gray-400" />
        </div>
        <input
          type="text"
          placeholder="Search..."
          className="block w-full rounded-md border-0 py-2 pl-10 text-gray-900 ring-1 ring-inset ring-gray-300 placeholder:text-gray-400 focus:ring-2 focus:ring-inset focus:ring-green-600"
        />
      </div>
      
      <div className="ml-auto flex items-center space-x-4">
        <button className="relative p-1 text-gray-400 hover:text-gray-500">
          <Bell className="h-6 w-6" />
          <span className="absolute top-0 right-0 h-2 w-2 rounded-full bg-red-500"></span>
        </button>
        
        <div className="flex items-center space-x-2">
          <div className="h-8 w-8 rounded-full bg-green-500 flex items-center justify-center text-white">
            <span className="font-semibold">A</span>
          </div>
          <span className="text-sm font-medium text-gray-700 hidden md:inline-block">Admin</span>
        </div>
      </div>
    </header>
  );
};

export default TopBar;
