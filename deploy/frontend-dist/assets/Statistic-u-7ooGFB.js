import{Ci as e,Cr as t,Or as n,Si as r,Sr as i,ca as a,ea as o,lr as s,or as c,sa as l}from"./router-B_0ApDLw.js";import{n as u}from"./index-DdmnAdeW.js";var d=r(`statistic`,[e(`label`,`
 font-weight: var(--n-label-font-weight);
 transition: .3s color var(--n-bezier);
 font-size: var(--n-label-font-size);
 color: var(--n-label-text-color);
 `),r(`statistic-value`,`
 margin-top: 4px;
 font-weight: var(--n-value-font-weight);
 `,[e(`prefix`,`
 margin: 0 4px 0 0;
 font-size: var(--n-value-font-size);
 transition: .3s color var(--n-bezier);
 color: var(--n-value-prefix-text-color);
 `,[r(`icon`,{verticalAlign:`-0.125em`})]),e(`content`,`
 font-size: var(--n-value-font-size);
 transition: .3s color var(--n-bezier);
 color: var(--n-value-text-color);
 `),e(`suffix`,`
 margin: 0 0 0 4px;
 font-size: var(--n-value-font-size);
 transition: .3s color var(--n-bezier);
 color: var(--n-value-suffix-text-color);
 `,[r(`icon`,{verticalAlign:`-0.125em`})])])]),f=l({name:`Statistic`,props:Object.assign(Object.assign({},c.props),{tabularNums:Boolean,label:String,value:[String,Number]}),slots:Object,setup(e){let{mergedClsPrefixRef:n,inlineThemeDisabled:r,mergedRtlRef:a}=t(e),l=c(`Statistic`,`-statistic`,d,u,e,n),f=s(`Statistic`,a,n),p=o(()=>{let{self:{labelFontWeight:e,valueFontSize:t,valueFontWeight:n,valuePrefixTextColor:r,labelTextColor:i,valueSuffixTextColor:a,valueTextColor:o,labelFontSize:s},common:{cubicBezierEaseInOut:c}}=l.value;return{"--n-bezier":c,"--n-label-font-size":s,"--n-label-font-weight":e,"--n-label-text-color":i,"--n-value-font-weight":n,"--n-value-font-size":t,"--n-value-prefix-text-color":r,"--n-value-suffix-text-color":a,"--n-value-text-color":o}}),m=r?i(`statistic`,void 0,p,e):void 0;return{rtlEnabled:f,mergedClsPrefix:n,cssVars:r?void 0:p,themeClass:m?.themeClass,onRender:m?.onRender}},render(){var e;let{mergedClsPrefix:t,$slots:{default:r,label:i,prefix:o,suffix:s}}=this;return(e=this.onRender)==null||e.call(this),a(`div`,{class:[`${t}-statistic`,this.themeClass,this.rtlEnabled&&`${t}-statistic--rtl`],style:this.cssVars},n(i,e=>a(`div`,{class:`${t}-statistic__label`},this.label||e)),a(`div`,{class:`${t}-statistic-value`,style:{fontVariantNumeric:this.tabularNums?`tabular-nums`:``}},n(o,e=>e&&a(`span`,{class:`${t}-statistic-value__prefix`},e)),this.value===void 0?n(r,e=>e&&a(`span`,{class:`${t}-statistic-value__content`},e)):a(`span`,{class:`${t}-statistic-value__content`},this.value),n(s,e=>e&&a(`span`,{class:`${t}-statistic-value__suffix`},e))))}});export{f as t};