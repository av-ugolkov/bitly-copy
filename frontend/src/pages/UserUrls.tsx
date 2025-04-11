import { DocumentDuplicateIcon, TrashIcon } from '@heroicons/react/24/outline';
import { useState } from 'react';

export default function UserUrls() {
  const [userID, setUserID] = useState('');
  const [links, setLinks] = useState<string[]>([]);

  async function getUserURLs() {
    setLinks([]);
    const resp = await fetch(`http://localhost:3000/get-urls?uid=${userID}`, {
      method: 'GET',
      headers: {
        'Content-Type': 'application/json',
        Accept: 'application/json',
      },
    });
    if (!resp.ok) {
      console.log(resp);
      return;
    }

    const data = await resp.json();
    data['urls'].forEach((link: string) => {
      setLinks((prev) => [...prev, link]);
    });
  }

  function copyToClipboard(text: string) {
    navigator.clipboard.writeText(text);
  }

  function deleteLink(link: string) {
    fetch('http://localhost:3000/remove', {
      method: 'DELETE',
      body: JSON.stringify({
        uid: userID,
        short_url: link,
      }),
      headers: {
        'Content-Type': 'application/json',
        Accept: 'application/json',
      },
    });
    setLinks((prev) => prev.filter((item) => item !== link));
  }

  return (
    <div className='flex flex-col w-1/3 h-fit min-w-[480px] my-4 p-4 gap-y-2 bg-amber-200 border-2 rounded-2xl'>
      <div className='font-bold'>User URLs</div>
      <input
        type='text'
        placeholder='Enter your nickname'
        onChange={(e) => setUserID(e.target.value)}
        className='w-full h-full bg-white p-1 border-0'
      />
      <button
        className='w-full h-full bg-white p-1 border-0'
        onClick={getUserURLs}>
        Get URLs
      </button>
      <label className='flex w-full h-full justify-between bg-gray-200 p-1 border-0'>
        <div className='w-full text-center'>
          {links.length > 0 ? (
            <div className='flex flex-col gap-y-2'>
              {links.map((link) => (
                <div className='flex w-full justify-between'>
                  <div className='flex'>{link}</div>
                  <div className='flex gap-x-2'>
                    <DocumentDuplicateIcon
                      className='size-6 text-gray-500 hover:text-green-500'
                      onClick={() => copyToClipboard(link)}
                    />
                    <TrashIcon
                      className='size-6 text-gray-500 hover:text-red-500'
                      onClick={() => deleteLink(link)}
                    />
                  </div>
                </div>
              ))}
            </div>
          ) : (
            'Here will be displayed your short links'
          )}
        </div>
      </label>
    </div>
  );
}
