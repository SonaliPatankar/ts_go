import React, { useState } from 'react';
import axios from 'axios';

type GitHubUser = {
  login: string;
  name: string;
  avatar_url: string;
  location: string;
  public_repos: number;
  followers: number;
  following: number;
};

function ApiCall() {
  const [username, setUsername] = useState('');
  const [user, setUser] = useState<GitHubUser | null>(null);

  const fetchUser = async () => {
    const res = await axios.get(`http://localhost:8080/github/user/${username}`);
    setUser(res.data);
  };

  return (
    <div style={{ padding: 20 }}>
      <h1>🔍 GitHub User Lookup</h1>
      <input
        value={username}
        onChange={(e) => setUsername(e.target.value)}
        placeholder="Enter GitHub username"
      />
      <button onClick={fetchUser}>Search</button>

      {user && (
        <div style={{ marginTop: 20 }}>
          <img src={user.avatar_url} alt="avatar" width={100} />
          <h2>{user.name} (@{user.login})</h2>
          <p>Location: {user.location}</p>
          <p>Repos: {user.public_repos} | Followers: {user.followers} | Following: {user.following}</p>
        </div>
      )}
    </div>
  );
}

export default ApiCall;
