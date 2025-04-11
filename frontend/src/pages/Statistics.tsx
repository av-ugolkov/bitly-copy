import { useState } from 'react';

export default function Statistics() {
  const [shortLink, setShortLink] = useState('');
  const [countRedirects, setCountRedirects] = useState(-1);

  async function getStatistics() {
    const resp = await fetch(
      `http://localhost:3000/statistics?short_url=${shortLink}`,
      {
        method: 'GET',
        headers: {
          'Content-Type': 'application/json',
          Accept: 'application/json',
        },
      }
    );
    if (!resp.ok) {
      console.log(resp);
      return;
    }

    const data = await resp.json();
    console.log(data);
    setCountRedirects(data['statistics']);
  }

  return (
    <div className='flex flex-col w-1/3 h-fit min-w-[480px] my-4 p-4 gap-y-2 bg-amber-200 border-2 rounded-2xl'>
      <div className='font-bold'>Statistics for the short link</div>
      <input
        type='text'
        placeholder='Enter the short link'
        onChange={(e) => {
          setCountRedirects(-1);
          setShortLink(e.target.value);
        }}
        className='w-full h-full bg-white p-1 border-0'
      />
      <button
        className='w-full h-full bg-white p-1 border-0'
        onClick={getStatistics}>
        Statistics
      </button>
      <label className='flex w-full h-full justify-between bg-gray-200 p-1 border-0'>
        <div className='w-full text-center'>
          {countRedirects != -1 ? (
            <div className='flex w-full justify-between'>
              <div className='flex'>{shortLink}</div>
              <div className='flex font-bold px-2'>{countRedirects}</div>
            </div>
          ) : (
            'Here will be displayed statistics for the short link'
          )}
        </div>
      </label>
    </div>
  );
}
