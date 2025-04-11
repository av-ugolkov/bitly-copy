import { DocumentDuplicateIcon } from '@heroicons/react/24/outline';
import { useState } from 'react';

export default function Main() {
  const [userID, setUserID] = useState('');
  const [longLink, setLongLink] = useState('');
  const [aliase, setAliase] = useState('');
  const [shortLink, setShortLink] = useState('');

  async function generateShortLink() {
    const resp = await fetch('http://localhost:3000/add', {
      method: 'POST',
      body: JSON.stringify({
        uid: userID,
        url: longLink,
        aliase: aliase,
      }),
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
    setShortLink(data['short_url']);
  }

  function copyToClipboard() {
    navigator.clipboard.writeText(shortLink);
  }

  return (
    <div className='flex flex-col w-1/3 h-fit min-w-[480px] my-4 p-4 gap-y-2 bg-amber-200 border-2 rounded-2xl'>
      <div className='font-bold'>Shorten a long URL</div>
      <input
        type='text'
        placeholder='Enter your nickname'
        onChange={(e) => setUserID(e.target.value)}
        className='w-full h-full bg-white p-1 border-0'
      />
      <input
        type='text'
        placeholder='Enter your long link'
        onChange={(e) => setLongLink(e.target.value)}
        className='w-full h-full bg-white p-1 border-0'
      />
      <input
        type='text'
        placeholder='Aliase (optional)'
        onChange={(e) => setAliase(e.target.value)}
        className='w-full h-full bg-white p-1 border-0'
      />
      <button
        className='w-full h-full bg-white p-1 border-0'
        onClick={generateShortLink}>
        Generate
      </button>
      <label className='flex w-full h-full justify-between bg-gray-200 p-1 border-0'>
        <div className='text-center'>
          {shortLink || 'Here will be displayed your short link'}
        </div>
        <DocumentDuplicateIcon
          className='size-6 text-gray-500 hover:text-green-500'
          onClick={copyToClipboard}
        />
      </label>
    </div>
  );
}
