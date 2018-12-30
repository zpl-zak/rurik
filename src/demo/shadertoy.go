package main

import (
	"github.com/zaklaus/rurik/src/core"
	"github.com/zaklaus/rurik/src/system"
)

type shadertoyProg struct {
	RenderTexture   system.RenderTarget
	ShadertoyShader system.Program
	frameIndex      int32
}

func newShadertoy() *shadertoyProg {
	return &shadertoyProg{
		RenderTexture:   system.CreateRenderTarget(320, 320),
		ShadertoyShader: system.NewProgramFromCode("", shadertoyHeaderCode+shadertoySrcCode),
	}
}

func (s *shadertoyProg) Apply() {
	s.ShadertoyShader.SetShaderValuei("iFrame", []int32{s.frameIndex}, 1)
	s.ShadertoyShader.RenderToTexture(core.NullTexture, s.RenderTexture)

	s.frameIndex++
}

const (
	shadertoyHeaderCode = `
#version 330
#extension GL_EXT_gpu_shader4 : enable

// Input vertex attributes (from vertex shader)
in vec2 fragTexCoord;
in vec4 fragColor;

// Input uniform values
uniform sampler2D texture0;
uniform vec4 colDiffuse;
uniform float time;
uniform int iFrame;

// Output fragment color
out vec4 finalColor;

uniform vec2 size = vec2(640, 480);

`
	shadertoySrcCode = `
// srtuss, 2013

#define ITER 20

vec2 circuit(vec2 p)
{
	p = fract(p);
	float r = 0.123;
	float v = 0.0, g = 0.0;
	float test = 0.0;
	r = fract(r * 9184.928);
	float cp, d;
	
	d = p.x;
	g += pow(clamp(1.0 - abs(d), 0.0, 1.0), 160.0);
	d = p.y;
	g += pow(clamp(1.0 - abs(d), 0.0, 1.0), 160.0);
	d = p.x - 1.0;
	g += pow(clamp(1.0 - abs(d), 0.0, 1.0), 160.0);
	d = p.y - 1.0;
	g += pow(clamp(1.0 - abs(d), 0.0, 1.0), 160.0);
	
	for(int i = 0; i < ITER; i ++)
	{
		cp = 0.5 + (r - 0.5) * 0.9;
		d = p.x - cp;
		g += pow(clamp(1.0 - abs(d), 0.0, 1.0), 160.0);
		if(d > 0.0)
		{
			r = fract(r * 4829.013);
			p.x = (p.x - cp) / (1.0 - cp);
			v += 1.0;
			test = r;
		}
		else
		{
			r = fract(r * 1239.528);
			p.x = p.x / cp;
			test = r;
		}
		p = p.yx;
	}
	
	v /= float(ITER);
	
	return vec2(v, g);
}

float rand12(vec2 p)
{
	vec2 r = (456.789 * sin(789.123 * p.xy));
	return fract(r.x * r.y);
}

vec2 rand22(vec2 p)
{
	vec2 ra = (456.789 * sin(789.123 * p.xy));
	vec2 rb = (456.789 * cos(789.123 * p.xy));
	return vec2(fract(ra.x * ra.y + p.x), fract(rb.x * rb.y + p.y));
}

float noise12(vec2 p)
{
	vec2 fl = floor(p);
	vec2 fr = fract(p);
	fr = smoothstep(0.0, 1.0, fr);
	return mix(
		mix(rand12(fl),                  rand12(fl + vec2(1.0, 0.0)), fr.x),
		mix(rand12(fl + vec2(0.0, 1.0)), rand12(fl + vec2(1.0, 1.0)), fr.x), fr.y);
}

float fbm12(vec2 p)
{
	return noise12(p) * 0.5 + noise12(p * 2.0 + vec2(11.93, 0.41)) * 0.25 + noise12(p * 4.0 + vec2(1.93, -17.41)) * 0.125 + noise12(p * 8.0 + vec2(-19.93, 11.41)) * 0.0625;
}
	
vec3 voronoi(in vec2 x)
{
	vec2 n = floor(x); // grid cell id
	vec2 f = fract(x); // grid internal position
	vec2 mg; // shortest distance...
	vec2 mr; // ..and second shortest distance
	float md = 8.0, md2 = 8.0;
	
	for(int j = -1; j <= 1; j ++)
	{
		for(int i = -1; i <= 1; i ++)
		{
			vec2 g = vec2(float(i), float(j)); // cell id
			vec2 o = rand22(n + g); // offset to edge point
			vec2 r = g + o - f;
			
			float d = max(abs(r.x), abs(r.y)); // distance to the edge
			
			if(d < md)
				{md2 = md; md = d; mr = r; mg = g;}
			else if(d < md2)
				{md2 = d;}
		}
	}
	return vec3(n + mg, md2 - md);
}

vec2 rotate(vec2 p, float a)
{
	return vec2(p.x * cos(a) - p.y * sin(a), p.x * sin(a) + p.y * cos(a));
}

float xor(float a, float b)
{
    return min(max(-a, b), max(a, -b));
}

float fr2(vec2 uv)
{
    float v = 1e38, dfscl = 1.0;
    
    vec4 rnd = vec4(0.1, 0.3, 0.7, 0.8);
    
    #define RNDA rnd = fract(sin(rnd * 11.1111) * 2986.3971)
    #define RNDB rnd = fract(cos(rnd * 11.1111) * 2986.3971)
    
    RNDA;
    
    for(int i = 0; i < 8; i++)
    {
        vec2 p = uv;
        
        //p.x += time;
        
        float si = 1.0 + rnd.x;
        p = (abs(fract(p / si) - 0.5)) * si;
        vec2 q = p;
        float w = max(q.x - rnd.y * 0.7, q.y - rnd.z * 0.7);
        w /= dfscl;
        v = xor(v, w);
        
        if(w < 0.0)
        {
        	RNDA;
        }
        else
        {
            RNDB;
        }
        
        float sii = 1.2;
        
        uv *= sii;
        uv -= rnd.xz;
        dfscl *= sii;
    }
    return v;
}


vec3 pixel(vec2 uv)
{
	uv.x *= size.x / size.y;
	
	float t = time * 0.5;
	
	uv = rotate(uv, sin(t) * 0.1);
	uv += t * vec2(0.5, 1.0);
	
	vec2 ci = circuit(uv * 0.1);
	
	vec3 vo, vo2, vo3;
	vo = voronoi(uv);
	
	float f = 80.0;
	
	float cf = 0.1;
	vec2 fr = (fract(uv / cf) - 0.5) * cf;
	vec2 fl = (floor(uv / cf) - 0.5) * cf;
	float cir = length(fr /*+ normalize(rand22(fl) * 2.0 - 1.0) * 0.2*/) - 0.03;
	
    float v;
    v = min(cos(vo.z * f), cir * 50.0) + ci.y;
	
    float ww = fr2(uv / 1.5) * 1.5;
    v = max(v, smoothstep(0.0, 0.01, ww - ci.y * 0.03));
    
    v = smoothstep(0.2, 0.0, v);
    
	//v = mix(v, 0.0, );
    //v = mix(v, 1.0, smoothstep(0.02, 0.0, abs(ww)));
	
	//v += smoothstep(0.01, 0.0, length(fr) - 0.1);
	
	return vec3(v);
}

void main()
{
	vec3 col;
	
	vec2 h = vec2(0.5, 0.0);
	col = pixel(fragTexCoord.xy + h.yy);
	col += pixel(fragTexCoord.xy + h.xy);
	col += pixel(fragTexCoord.xy + h.yx);
	col += pixel(fragTexCoord.xy + h.xx);
	
	col /= 4.0;
	
	//col = vec3(1.0);
	
	vec2 uv = fragTexCoord.xy;
	col *= ((1.0 - pow(abs(uv.x), 2.1)) * (1.0 - pow(abs(uv.y), 2.1)));
	
	finalColor = vec4(col, 1.0);
}
	`
)
